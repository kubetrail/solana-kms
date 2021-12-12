package run

import (
	"bytes"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	kms "cloud.google.com/go/kms/apiv1"
	"github.com/kubetrail/solana-kms/pkg/flags"
	"github.com/portto/solana-go-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	kms2 "google.golang.org/genproto/googleapis/cloud/kms/v1"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// KeyNew generates a new private keypair data either from random seed or a seedfile
// provided as input. The seedfile needs to be in encrypted format.
// When seedfile is not provided, a seedfile is generated along with the private keypair
// data with .seed extension
func KeyNew(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	persistentFlags := getPersistentFlags(cmd)

	_ = viper.BindPFlag(flags.KeyFile, cmd.Flags().Lookup(filepath.Base(flags.KeyFile)))
	_ = viper.BindPFlag(flags.SeedFile, cmd.Flags().Lookup(filepath.Base(flags.SeedFile)))

	configFile := viper.GetString(flags.Config)
	keyFile := viper.GetString(flags.KeyFile)
	seedFile := viper.GetString(flags.SeedFile)

	if err := setAppCredsEnvVar(persistentFlags.ApplicationCredentials); err != nil {
		err := fmt.Errorf("could not set Google Application credentials env. var: %w", err)
		return err
	}

	if len(keyFile) == 0 {
		if len(configFile) == 0 {
			var err error
			configFile, err = getDefaultConfigFilename()
			if err != nil {
				err := fmt.Errorf("could not get default config filename: %w", err)
				return err
			}
		}

		configValues, err := getConfigValues(configFile)
		if err != nil {
			err := fmt.Errorf("could not get config values: %w", err)
			return err
		}

		if configValues == nil || len(configValues.KeypairPath) == 0 {
			err := fmt.Errorf("could not find a valid keypair path from config file")
			return err
		}

		keyFile = configValues.KeypairPath
	}

	keyFile = removeSchemeFromPath(keyFile)

	kmsClient, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		err := fmt.Errorf("failed to create kms client: %w", err)
		return err
	}
	defer kmsClient.Close()

	var account types.Account

	// if seed file is provided, use that to generate new account
	// otherwise generate new account using default random seed
	if len(seedFile) > 0 {
		ciphertext, err := os.ReadFile(seedFile)
		if err != nil {
			err := fmt.Errorf("could not read seed file: %w", err)
			return err
		}

		decryptResponse, err := kmsClient.Decrypt(
			ctx,
			&kms2.DecryptRequest{
				Name: getKmsName(
					persistentFlags.Project,
					persistentFlags.Location,
					persistentFlags.Keyring,
					persistentFlags.Key,
				),
				Ciphertext:                        ciphertext,
				AdditionalAuthenticatedData:       nil,
				CiphertextCrc32C:                  wrapperspb.Int64(int64(crc32Sum(ciphertext))),
				AdditionalAuthenticatedDataCrc32C: nil,
			},
		)
		if err != nil {
			err := fmt.Errorf("could not decrypt seed: %w", err)
			return err
		}

		_, X, err := ed25519.GenerateKey(bytes.NewReader(decryptResponse.Plaintext))
		if err != nil {
			err := fmt.Errorf("could not generate ed25519 key: %w", err)
			return err
		}
		account, err = types.AccountFromBytes(X)
		if err != nil {
			err := fmt.Errorf("could not generate new account: %w", err)
			return err
		}
	} else {
		account = types.NewAccount()
	}

	encryptResponseKey, err := kmsClient.Encrypt(
		ctx,
		&kms2.EncryptRequest{
			Name: getKmsName(
				persistentFlags.Project,
				persistentFlags.Location,
				persistentFlags.Keyring,
				persistentFlags.Key,
			),
			Plaintext:                         account.PrivateKey,
			AdditionalAuthenticatedData:       nil,
			PlaintextCrc32C:                   wrapperspb.Int64(int64(crc32Sum(account.PrivateKey))),
			AdditionalAuthenticatedDataCrc32C: nil,
		},
	)
	if err != nil {
		err := fmt.Errorf("could not encrypt private key: %w", err)
		return err
	}

	encryptResponseSeed, err := kmsClient.Encrypt(
		ctx,
		&kms2.EncryptRequest{
			Name: getKmsName(
				persistentFlags.Project,
				persistentFlags.Location,
				persistentFlags.Keyring,
				persistentFlags.Key,
			),
			Plaintext:                         account.PrivateKey.Seed(),
			AdditionalAuthenticatedData:       nil,
			PlaintextCrc32C:                   wrapperspb.Int64(int64(crc32Sum(account.PrivateKey.Seed()))),
			AdditionalAuthenticatedDataCrc32C: nil,
		},
	)
	if err != nil {
		err := fmt.Errorf("could not encrypt private key: %w", err)
		return err
	}

	if keyFile == "-" {
		type keyInfo struct {
			PrivateKeyCipherText []byte `json:"privateKeyCipherText,omitempty"`
			SeedCipherText       []byte `json:"seedCipherText,omitempty"`
		}

		info := &keyInfo{
			PrivateKeyCipherText: encryptResponseKey.Ciphertext,
			SeedCipherText:       encryptResponseSeed.Ciphertext,
		}

		jb, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			err := fmt.Errorf("could not serialize key info: %w", err)
			return err
		}

		if _, err := fmt.Fprintln(
			cmd.OutOrStdout(),
			string(jb),
		); err != nil {
			err := fmt.Errorf("could not print base64 encoded private key: %w", err)
			return err
		}

		return nil
	}

	if err := os.WriteFile(keyFile, encryptResponseKey.Ciphertext, 0400); err != nil {
		err := fmt.Errorf("could not write encrypted private key to outfile: %w", err)
		return err
	}

	seedFile = fmt.Sprintf("%s.%s", keyFile, "seed")
	if err := os.WriteFile(seedFile, encryptResponseSeed.Ciphertext, 0400); err != nil {
		err := fmt.Errorf("could not write encrypted private key seed to outfile: %w", err)
		return err
	}

	return nil
}
