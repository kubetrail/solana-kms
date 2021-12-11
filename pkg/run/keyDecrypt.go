package run

import (
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

// KeyDecrypt decrypts KMS encrypted private keypair file and prints on screen the values
// of public key and the private key.
func KeyDecrypt(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	persistentFlags := getPersistentFlags(cmd)

	_ = viper.BindPFlag(flags.KeyFile, cmd.Flags().Lookup(filepath.Base(flags.KeyFile)))

	configFile := viper.GetString(flags.Config)
	keyFile := viper.GetString(flags.KeyFile)

	if err := setAppCredsEnvVar(persistentFlags.ApplicationCredentials); err != nil {
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
		err := fmt.Errorf("could not create key pair from decrypted data: %w", err)
		return err
	}

	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Public key: %s\n", account.PublicKey); err != nil {
		err := fmt.Errorf("could not print public key to cmd stdout: %w", err)
		return err
	}

	keyValues := make([]int, len(account.PrivateKey))
	for i, value := range account.PrivateKey {
		keyValues[i] = int(value)
	}

	jb, err := json.Marshal(keyValues)
	if err != nil {
		err := fmt.Errorf("could not serialize private key for displaying as array: %w", err)
		return err
	}

	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Private key: %s\n", string(jb)); err != nil {
		err := fmt.Errorf("could not print private key to cmd stdout: %w", err)
		return err
	}

	return nil
}
