package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.BoolFlag{Name: "debug", Usage: "debug commands"},
	}
	app.Commands = []cli.Command{
		{
			Name:    "create",
			Aliases: []string{"c", "add"},
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "c, current", Usage: "use current path to add"},
				cli.BoolFlag{Name: "e, editor", Usage: "use default editor to add"},
				cli.BoolFlag{Name: "v, validate-path", Usage: "this option allow to validate if path exists"},
			},
			Usage:     "create new project",
			UsageText: "project create <name> <path>",
			ArgsUsage: "name path",
			Action:    create,
		},
		{
			Name:    "update",
			Aliases: []string{"u"},
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "v, validate-path", Usage: "this option allow to validate if path exists"},
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
				cli.BoolFlag{Name: "f, full", Usage: "show full info about project"},
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
				cli.BoolFlag{Name: "d, duplicate", Usage: "duplicate running session"},
				cli.BoolFlag{Name: "vim", Usage: "open tmux with vim opened"},
			},
			Usage:     "open project using tmux",
			UsageText: "project open <name>",
			ArgsUsage: "name",
			Action:    open,
		},
		{
			Name:    "close",
			Aliases: []string{"x", "deattach"},
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "a, all", Usage: "all opened projects"},
				cli.BoolFlag{Name: "k, kill", Usage: "kill running project"},
			},
			Usage:     "close project",
			UsageText: "close project <?name>",
			ArgsUsage: "name",
			Action:    close,
		},
		{ //TODO essa opção editar não ficou legal, o editar deveria editar o projeto como é feito no update
			//TODO para abrir no vscode deveria ter alguma coisa diferente
			Name:      "edit",
			Aliases:   []string{"e"},
			Usage:     "edit project using vscode",
			UsageText: "project edit <name>",
			ArgsUsage: "name",
			Action:    edit,
		},
		{
			Name: "export",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "f, format", Usage: "export projects to (nerdtree|vimcommand)"},
				cli.BoolFlag{Name: "override", Usage: "Override default configuration file"},
			},
			Usage:     "export projects to use in another locations",
			UsageText: "project -f <format>",
			Action:    export,
		},
		{
			Name:      "path",
			Aliases:   []string{"pt"},
			Usage:     "show path of project",
			UsageText: "project path <name>",
			ArgsUsage: "name",
			Action:    path,
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
