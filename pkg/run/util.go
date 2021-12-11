package run

import (
	"fmt"
	"hash/crc32"
	"os"
	"path/filepath"
	"strings"

	"github.com/portto/solana-go-sdk/rpc"

	"github.com/kubetrail/solana-kms/pkg/flags"
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
	ApplicationCredentials string `json:"applicationCredentials,omitempty"`
	Project                string `json:"project,omitempty"`
	Location               string `json:"location,omitempty"`
	Keyring                string `json:"keyring,omitempty"`
	Key                    string `json:"key,omitempty"`
}

func getPersistentFlags(cmd *cobra.Command) persistentFlagValues {
	_ = viper.BindPFlag(flags.Config, cmd.PersistentFlags().Lookup(filepath.Base(flags.Config)))
	_ = viper.BindPFlag(flags.Project, cmd.PersistentFlags().Lookup(filepath.Base(flags.Project)))
	_ = viper.BindPFlag(flags.Location, cmd.PersistentFlags().Lookup(filepath.Base(flags.Location)))
	_ = viper.BindPFlag(flags.Keyring, cmd.PersistentFlags().Lookup(filepath.Base(flags.Keyring)))
	_ = viper.BindPFlag(flags.Key, cmd.PersistentFlags().Lookup(filepath.Base(flags.Key)))
	_ = viper.BindPFlag(flags.ApplicationCredentials, cmd.PersistentFlags().Lookup(filepath.Base(flags.ApplicationCredentials)))

	_ = viper.BindEnv(flags.Project, "PROJECT")
	_ = viper.BindEnv(flags.Location, "LOCATION")
	_ = viper.BindEnv(flags.Keyring, "KEYRING")
	_ = viper.BindEnv(flags.Key, "KEY")

	applicationCredentials := viper.GetString(flags.ApplicationCredentials)
	project := viper.GetString(flags.Project)
	location := viper.GetString(flags.Location)
	keyring := viper.GetString(flags.Keyring)
	key := viper.GetString(flags.Key)

	return persistentFlagValues{
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
	case "mainnet":
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
