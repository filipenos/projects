package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(codeCmd)
	codeCmd.Flags().StringP("editor", "e", "code", "Code using vscode editor")
	codeCmd.Flags().BoolP("reuse", "r", false, "Reuse same window")
}

// codeCmd represents the code command
var codeCmd = &cobra.Command{
	Use:   "code",
	Short: "Edit your project using the editor (code as default)",
	RunE:  code,
}

func code(cmdParam *cobra.Command, params []string) error {
	projects, err := Load(LoadSettings())
	if err != nil {
		return err
	}

	p, _ := projects.Find(safeName(params...))
	if p == nil {
		return errorf("project not found")

	}
	if p.Path == "" {
		return errorf("project '%s' dont have path", p.Name)
	}
	if !isExist(p.Path) {
		return errorf("path '%s' of project '%s' not exists", p.Path, p.Name)
	}

	editor := "code"
	args := make([]string, 0)

	if edit := SafeStringFlag(cmdParam, "e"); edit != "" {
		editor = edit
	} else {
		pos := "--new-window"
		if SafeBoolFlag(cmdParam, "r") {
			pos = "--reuse-window"
		}
		args = append(args, pos)
	}
	args = append(args, p.Path)

	log("open path '%s' on '%s'", p.Path, editor)

	cmd := exec.Command(editor, args...)
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
