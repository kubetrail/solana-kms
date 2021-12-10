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

// keyPubkeyCmd represents the keyPubkey command
var keyPubkeyCmd = &cobra.Command{
	Use:   "pubkey",
	Short: "Get public key from KMS encrypted keypair",
	Long: `This command fetches public key associated with the private key.

Assuming you have generated a keypair at /tmp/key, the public key can be obtained
as follows:

solana-kms key pubkey --keyfile=/tmp/key
`,
	RunE: run.KeyPubkey,
}

func init() {
	keyCmd.AddCommand(keyPubkeyCmd)
	f := keyPubkeyCmd.Flags()
	b := filepath.Base

	f.String(b(flags.KeyFile), "", "Input key file")
}
