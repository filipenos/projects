/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generates bash completion scripts",
	Long: `To load completion run

. <(projects completion)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(projects completion)
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		sh, filename := "zsh", SafeStringFlag(cmd, "file")
		if len(args) > 0 {
			sh = args[0]
		}

		switch sh {
		case "bash":
			if filename == "" {
				return rootCmd.GenBashCompletion(os.Stdout)
			}
			return rootCmd.GenBashCompletionFile(filename)

		case "fish":
			if filename == "" {
				return rootCmd.GenFishCompletion(os.Stdout, true)
			}
			return rootCmd.GenFishCompletionFile(filename, true)

		case "powershell":
			if filename == "" {
				return rootCmd.GenPowerShellCompletion(os.Stdout)
			}
			return rootCmd.GenPowerShellCompletionFile(filename)

		case "zsh":
			if filename == "" {
				return rootCmd.GenZshCompletion(os.Stdout)
			}
			return rootCmd.GenZshCompletionFile(filename)

		default:
			return fmt.Errorf("invalid shell %s", sh)
		}
	},
}

func init() {
	completionCmd.Flags().StringP("file", "f", "", "Generate completion on file")

	rootCmd.AddCommand(completionCmd)
}
