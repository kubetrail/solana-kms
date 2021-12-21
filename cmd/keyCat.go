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

// keyCatCmd represents the keyCat command
var keyCatCmd = &cobra.Command{
	Use:   "cat",
	Short: "Print contents of key",
	Long: `This command simply prints contents of the key irrespective
of whether it is encrypted or not.`,
	RunE: run.KeyCat,
}

func init() {
	keyCmd.AddCommand(keyCatCmd)
	f := keyCatCmd.Flags()
	b := filepath.Base

	f.String(b(flags.KeyFile), "", "Input key file")
}
