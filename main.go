package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

//TODO adicionar suporte a outra extensao de projetdos do vscode

//TODO pensar em como lidar com projetos desabilitados.
//ele não é listado, mas é necessário um lugar para remover, ou habilitar novamente
//o que está acontecendo agora, é que se o parametro existir no json ele é utilizado
//talvez permitir de editar o arquivo json já seria suficiente.
//pra isso a opção edit talvez
func main() {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		&cli.BoolFlag{Name: "debug", Usage: "debug commands"},
	}
	app.Commands = []*cli.Command{
		{
			Name:    "create",
			Aliases: []string{"c", "add"},
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "e, editor", Usage: "use default editor to add"},
				&cli.BoolFlag{Name: "n, no-validate", Usage: "this option ignore path validation"},
			},
			Usage:     "create new project",
			UsageText: "project create <name> <path>",
			ArgsUsage: "name path",
			Action:    create,
		},
		{
			Name:    "update",
			Aliases: []string{"u", "edit", "e"},
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "n, no-validate", Usage: "this option ignore path validation"},
			},
			Usage:     "update project",
			UsageText: "project update <name>",
			ArgsUsage: "name",
			Action:    update,
		},
		{
			Name:      "delete",
			Aliases:   []string{"d", "rm"},
			Usage:     "delete project",
			UsageText: "project delete <name>",
			ArgsUsage: "name",
			Action:    delete,
		},
		{
			Name:    "list",
			Aliases: []string{"l", "ls"},
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "s, simple", Usage: "show only name of project"},
				&cli.BoolFlag{Name: "p, path", Usage: "show path of project"},
			},
			Usage:     "list projects",
			UsageText: "project list <options>",
			ArgsUsage: "<options>",
			Action:    list,
		},
		{
			Name:    "open",
			Aliases: []string{"o", "attach"},
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "d, duplicate", Usage: "duplicate running session"},
				&cli.BoolFlag{Name: "r, recursively", Usage: "if name is empty, search recursively on path"},
			},
			Usage:     "open project using tmux",
			UsageText: "project open <name>",
			ArgsUsage: "name",
			Action:    open,
		},
		{
			Name:    "code",
			Aliases: []string{"vscode"},
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "r, recursively", Usage: "if name is empty, search recursively on path"},
				&cli.StringFlag{Name: "e, editor", Usage: "change default editor (vscode)"},
			},
			Usage:     "open project using code",
			UsageText: "project code <name>",
			ArgsUsage: "name",
			Action:    code,
		},
		{
			Name:    "close",
			Aliases: []string{"x", "deattach"},
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "a, all", Usage: "all opened projects"},
				&cli.BoolFlag{Name: "k, kill", Usage: "kill running project"},
			},
			Usage:     "close project",
			UsageText: "close project <name>",
			ArgsUsage: "name",
			Action:    close,
		},
		{
			Name: "export",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "format", Usage: "export projects to (nerdtree|vimcommand|vim-project|aliases)"},
				&cli.BoolFlag{Name: "override", Usage: "Override default configuration file"},
			},
			Usage:     "export projects to use in another locations",
			UsageText: "project -f <format>",
			Action:    export,
		},
		{
			Name:      "path",
			Aliases:   []string{"pt", "pwd"},
			Usage:     "show path of project",
			UsageText: "project path <name>",
			ArgsUsage: "name",
			Action:    path,
		},
		{
			Name:    "scm",
			Aliases: []string{"git"},
			Flags: []cli.Flag{
				&cli.BoolFlag{Name: "set", Usage: ""},
			},
			Action: scm,
		},
	}
	app.Name = "Projects"
	app.Usage = "Simple manager for your projects"
	app.Description = "Manage local projects"
	app.HideVersion = true
	app.ExitErrHandler = func(c *cli.Context, err error) {
		if err != nil {
			log("%v", err)
		}
		return
	}
	app.Run(os.Args)
}
