package run

import (
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

// KeyPubkey fetches public key from the encrypted private keypair file
func KeyPubkey(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	persistentFlags := getPersistentFlags(cmd)

	_ = viper.BindPFlag(flags.KeyFile, cmd.Flags().Lookup(filepath.Base(flags.KeyFile)))

	keyFile := viper.GetString(flags.KeyFile)

	if err := setAppCredsEnvVar(persistentFlags.ApplicationCredentials); err != nil {
		err := fmt.Errorf("could not set Google Application credentials env. var: %w", err)
		return err
	}

	if len(keyFile) == 0 {
		if len(persistentFlags.ConfigFile) == 0 {
			var err error
			persistentFlags.ConfigFile, err = getDefaultConfigFilename()
			if err != nil {
				err := fmt.Errorf("could not get default config filename: %w", err)
				return err
			}
		}

		configValues, err := getConfigValues(persistentFlags.ConfigFile)
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

	ciphertext, err := os.ReadFile(keyFile)
	if err != nil {
		err := fmt.Errorf("error reading input keypair file: %w", err)
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
		err := fmt.Errorf("could not decrypt private key: %w", err)
		return err
	}

	account, err := types.AccountFromBytes(decryptResponse.Plaintext)
	if err != nil {
		err := fmt.Errorf("could not create account from data: %w", err)
		return err
	}

	if _, err := fmt.Fprintln(cmd.OutOrStdout(), account.PublicKey.ToBase58()); err != nil {
		err := fmt.Errorf("could not write to cmd output: %w", err)
		return err
	}

	return nil
}
