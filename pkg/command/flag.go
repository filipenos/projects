package command

import "github.com/spf13/cobra"

func SafeBoolFlag(cmd *cobra.Command, flagName string) bool {
	v, _ := cmd.Flags().GetBool(flagName)
	return v
}

func SafeStringFlag(cmd *cobra.Command, flagName string) string {
	v, _ := cmd.Flags().GetString(flagName)
	return v
}
