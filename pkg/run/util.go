package run

import (
	"fmt"
	"hash/crc32"
	"os"
	"path/filepath"
	"strings"

	"github.com/kubetrail/solana-kms/pkg/flags"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sigs.k8s.io/yaml"
)

// getConfigValues reads config file and returns a data structure
func getConfigValues(configFile string) (*config, error) {
	cc := &config{}

	b, err := os.ReadFile(configFile)
	if err != nil {
		err := fmt.Errorf("could not read solana config file: %w", err)
		return nil, err
	}

	if err := yaml.Unmarshal(b, cc); err != nil {
		err := fmt.Errorf("could not unmarshal solana config file: %w", err)
		return nil, err
	}

	return cc, nil
}

// getDefaultConfigFilename retrieves default config filename
func getDefaultConfigFilename() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		err := fmt.Errorf("could not get user home dir: %w", err)
		return "", err
	}

	return filepath.Join(homeDir, ".config", "solana", "cli", "config.yml"), nil
}

// crc32Sum produces crc32 sum
func crc32Sum(data []byte) uint32 {
	t := crc32.MakeTable(crc32.Castagnoli)
	return crc32.Checksum(data, t)
}

// getKmsName constructs the canonical URI endpoint path for KMS encryption call
func getKmsName(projectId, kmsLocation, keyringName, keyName string) string {
	return fmt.Sprintf(
		"projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
		projectId,
		kmsLocation,
		keyringName,
		keyName,
	)
}

type persistentFlagValues struct {
	ConfigFile             string `json:"configFile,omitempty"`
	ApplicationCredentials string `json:"applicationCredentials,omitempty"`
	Project                string `json:"project,omitempty"`
	Location               string `json:"location,omitempty"`
	Keyring                string `json:"keyring,omitempty"`
	Key                    string `json:"key,omitempty"`
}

func getPersistentFlags(cmd *cobra.Command) persistentFlagValues {
	rootCmd := cmd.Root().PersistentFlags()
	b := filepath.Base

	_ = viper.BindPFlag(flags.Config, rootCmd.Lookup(b(flags.Config)))
	_ = viper.BindPFlag(flags.GoogleProjectID, rootCmd.Lookup(b(flags.GoogleProjectID)))
	_ = viper.BindPFlag(flags.KmsLocation, rootCmd.Lookup(b(flags.KmsLocation)))
	_ = viper.BindPFlag(flags.KmsKeyring, rootCmd.Lookup(b(flags.KmsKeyring)))
	_ = viper.BindPFlag(flags.KmsKey, rootCmd.Lookup(b(flags.KmsKey)))
	_ = viper.BindPFlag(flags.GoogleApplicationCredentials, rootCmd.Lookup(b(flags.GoogleApplicationCredentials)))

	_ = viper.BindEnv(flags.Config, "SOLANA_CONFIG")
	_ = viper.BindEnv(flags.GoogleProjectID, "GOOGLE_PROJECT_ID")
	_ = viper.BindEnv(flags.KmsLocation, "KMS_LOCATION")
	_ = viper.BindEnv(flags.KmsKeyring, "KMS_KEYRING")
	_ = viper.BindEnv(flags.KmsKey, "KMS_KEY")

	configFile := viper.GetString(flags.Config)
	applicationCredentials := viper.GetString(flags.GoogleApplicationCredentials)
	project := viper.GetString(flags.GoogleProjectID)
	location := viper.GetString(flags.KmsLocation)
	keyring := viper.GetString(flags.KmsKeyring)
	key := viper.GetString(flags.KmsKey)

	return persistentFlagValues{
		ConfigFile:             configFile,
		ApplicationCredentials: applicationCredentials,
		Project:                project,
		Location:               location,
		Keyring:                keyring,
		Key:                    key,
	}
}

func setAppCredsEnvVar(applicationCredentials string) error {
	if len(applicationCredentials) > 0 {
		if err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", applicationCredentials); err != nil {
			err := fmt.Errorf("could not set Google Application credentials env. var: %w", err)
			return err
		}
	}

	return nil
}

func getEndpointFromUrlOrMoniker(url string, configValues *config) string {
	switch strings.ToLower(url) {
	case "mainnet", "mainnet-beta":
		return rpc.MainnetRPCEndpoint
	case "devnet":
		return rpc.DevnetRPCEndpoint
	case "testnet":
		return rpc.TestnetRPCEndpoint
	case "localnet", "localhost":
		return rpc.LocalnetRPCEndpoint
	case "":
		return configValues.JsonRpcUrl
	default:
		return url
	}
}

func removeSchemeFromPath(input string) string {
	return strings.TrimLeft(input, "stdin:")
}
