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

// keyNewCmd represents the keyNew command
var keyNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new KMS encrypted keypair",
	Long: `This command will generate a new keypair and encrypt
using Google KMS. The persisted file is ciphertext.

First set following environment variables:
export GOOGLE_APPLICATION_CREDENTIALS="kms-encrypter-decrypter-role-sa.json"
export LOCATION=<kms location>
export KEYRING=<keyring name>
export KEY=<key name>
export PROJECT=<project id>

Generate a new key as follows
solana-kms key new --keyfile=/tmp/key

This will write two files: /tmp/key and /tmp/key.seed and both are encrypted using KMS.

The seed file can be used to recover the key later as follows:
solana-kms key new --keyfile=/tmp/key-recovered --seedfile=/tmp/key.seed

Please note that since each of these files is encrypted the contents of the original
and recovered files will not match. Decrypting both will, however, result in 
exact same contents.

You can verify as follows:
gcloud kms decrypt \
	--location=<location> \
	--keyring=<keyring name> \
	--key=<key name> \
	--ciphertext-file=/tmp/key \
	--plaintext-file=- \
| base64

gcloud kms decrypt \
	--location=<location> \
	--keyring=<keyring name> \
	--key=<key name> \
	--ciphertext-file=/tmp/key-recovered \
	--plaintext-file=- \
| base64
`,
	RunE: run.KeyNew,
}

func init() {
	keyCmd.AddCommand(keyNewCmd)
	f := keyNewCmd.Flags()
	b := filepath.Base

	f.String(b(flags.KeyFile), "", "Output key file")
	f.String(b(flags.SeedFile), "", "Input seed file")
}
