package run

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kubetrail/solana-kms/pkg/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// KeyCat prints contents of the key whether it is encrypted or not
func KeyCat(cmd *cobra.Command, _ []string) error {
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

	b, err := os.ReadFile(keyFile)
	if err != nil {
		err := fmt.Errorf("error reading input keypair file: %w", err)
		return err
	}

	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\n", string(b)); err != nil {
		err := fmt.Errorf("could not print key to cmd stdout: %w", err)
		return err
	}

	return nil
}
