package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
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
			Name:      "open",
			Aliases:   []string{"o"},
			Usage:     "open project",
			UsageText: "project open <name>",
			ArgsUsage: "name",
			Action:    open,
		},
		{
			Name:      "edit",
			Aliases:   []string{"e"},
			Usage:     "edit project",
			UsageText: "project edit <name>",
			ArgsUsage: "name",
			Action:    edit,
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
