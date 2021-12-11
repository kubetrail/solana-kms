package run

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	kms "cloud.google.com/go/kms/apiv1"
	"github.com/kubetrail/solana-kms/pkg/flags"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	kms2 "google.golang.org/genproto/googleapis/cloud/kms/v1"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// AccountInfo retrieves account info
func AccountInfo(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()
	persistentFlags := getPersistentFlags(cmd)

	_ = viper.BindPFlag(flags.KeyFile, cmd.Flags().Lookup(filepath.Base(flags.KeyFile)))
	_ = viper.BindPFlag(flags.PubKey, cmd.Flags().Lookup(filepath.Base(flags.PubKey)))
	_ = viper.BindPFlag(flags.Url, cmd.Flags().Lookup(filepath.Base(flags.Url)))

	configFile := viper.GetString(flags.Config)
	keyFile := viper.GetString(flags.KeyFile)
	pubKey := viper.GetString(flags.PubKey)
	url := viper.GetString(flags.Url)

	var endpoint string
	var configValues *config

	if len(pubKey) == 0 || len(url) == 0 {
		var err error
		if len(configFile) == 0 {
			configFile, err = getDefaultConfigFilename()
			if err != nil {
				err := fmt.Errorf("could not get default config filename: %w", err)
				return err
			}
		}

		configValues, err = getConfigValues(configFile)
		if err != nil {
			err := fmt.Errorf("could not get config values: %w", err)
			return err
		}
	}

	switch strings.ToLower(url) {
	case "mainnet":
		endpoint = rpc.MainnetRPCEndpoint
	case "devnet":
		endpoint = rpc.DevnetRPCEndpoint
	case "testnet":
		endpoint = rpc.TestnetRPCEndpoint
	case "localnet", "localhost":
		endpoint = rpc.LocalnetRPCEndpoint
	case "":
		endpoint = configValues.JsonRpcUrl
	default:
		endpoint = url
	}

	if len(pubKey) == 0 {
		if err := setAppCredsEnvVar(persistentFlags.ApplicationCredentials); err != nil {
			err := fmt.Errorf("could not set Google Application credentials env. var: %w", err)
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
			err := fmt.Errorf("could not create account from data: %w", err)
			return err
		}

		pubKey = account.PublicKey.ToBase58()
	}

	// create a RPC client
	c := client.NewClient(endpoint)

	response, err := c.GetAccountInfo(ctx, pubKey)
	if err != nil {
		err := fmt.Errorf("could not get account balance: %w", err)
		return err
	}

	jb, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		err := fmt.Errorf("could not serialize account info: %w", err)
		return err
	}

	if _, err := fmt.Fprintln(cmd.OutOrStdout(), string(jb)); err != nil {
		err := fmt.Errorf("could not write to command out: %w", err)
		return err
	}

	return nil
}
