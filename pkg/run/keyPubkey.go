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
	bindPersistentFlags(cmd)

	_ = viper.BindPFlag(flags.KeyFile, cmd.Flags().Lookup(filepath.Base(flags.KeyFile)))

	configFile := viper.GetString(flags.Config)
	keyFile := viper.GetString(flags.KeyFile)

	applicationCredentials := viper.GetString(flags.ApplicationCredentials)
	project := viper.GetString(flags.Project)
	location := viper.GetString(flags.Location)
	keyring := viper.GetString(flags.Keyring)
	key := viper.GetString(flags.Key)

	if err := setAppCredsEnvVar(applicationCredentials); err != nil {
		err := fmt.Errorf("could not set Google Application credentials env. var: %w", err)
		return err
	}

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

	if len(keyFile) == 0 {
		keyFile = configValues.KeypairPath
	}

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
				project,
				location,
				keyring,
				key,
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
