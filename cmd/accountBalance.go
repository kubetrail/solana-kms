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
	"github.com/kubetrail/solana-kms/pkg/run"
	"github.com/spf13/cobra"
)

// accountBalanceCmd represents the accountBalance command
var accountBalanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Find account balance",
	Long:  `This commands finds account balance`,
	RunE:  run.AccountBalance,
}

func init() {
	accountCmd.AddCommand(accountBalanceCmd)
	f := accountBalanceCmd.Flags()
	b := filepath.Base

	f.String(b(flags.KeyFile), "", "Keypair file")
	f.String(b(flags.PubKey), "", "Public key (--keyfile will be ignored)")
	f.String(b(flags.Url), "", "Solana validator endpoint (--keyfile will be ignored)")
}
