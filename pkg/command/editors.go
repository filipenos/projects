package command

import (
	"github.com/filipenos/projects/pkg/log"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "editors",
		Short: "List supported editors and their availability",
		RunE: func(cmd *cobra.Command, args []string) error {
			available, notAvailable := editorService.GetEditors()

			if len(available) > 0 {
				log.Println("Available:")
				for _, name := range available {
					log.Printf("  %-15s (found)\n", name)
				}
			}

			if len(notAvailable) > 0 {
				log.Println("Not available:")
				for _, name := range notAvailable {
					log.Printf("  %-15s (not found in PATH)\n", name)
				}
			}

			return nil
		},
	}
	rootCmd.AddCommand(cmd)
}
