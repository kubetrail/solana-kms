/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"path/filepath"

	"github.com/kubetrail/solana-kms/pkg/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "solana-kms",
	Short: "cli to work with solana network with Google KMS integration",
	Long: `This is a CLI tool that ensures that private keys are never
written to the disk. Key generation happens in memory and is encrypted
via Google KMS. All subsequent actions assume persisted keypair file to
be in the ciphertext format.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	f := rootCmd.PersistentFlags()
	b := filepath.Base

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	f.String(b(flags.GoogleProjectID), "", "Google project ID (Env: GOOGLE_PROJECT_ID)")
	f.String(b(flags.GoogleApplicationCredentials), "", "Google app credentials (Env: GOOGLE_APPLICATION_CREDENTIALS)")
	f.String(b(flags.KmsLocation), "global", "KMS location (Env: KMS_LOCATION)")
	f.String(b(flags.KmsKeyring), "", "KMS keyring name (Env: KMS_KEYRING)")
	f.String(b(flags.KmsKey), "", "KMS key name (Env: KMS_KEY)")

	f.String(b(flags.Config), "", "Solana config file (Env: SOLANA_CONFIG)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	viper.AutomaticEnv() // read in environment variables that match
}
