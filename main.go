package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

var (
	filepath = fmt.Sprintf("%s/.projects.json", os.Getenv("HOME"))
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:    "add",
			Aliases: []string{"a"},
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "current", Usage: "use current path to add"},
			},
			Usage:     "add new project",
			UsageText: "project add <name> <path>",
			ArgsUsage: "name path",
			Action:    add,
		},
		{
			Name:      "edit",
			Aliases:   []string{"e"},
			Usage:     "edit project config",
			UsageText: "project edit <name>",
			ArgsUsage: "name",
			Action:    edit,
		},
		{
			Name:      "remove",
			Aliases:   []string{"r"},
			Usage:     "remove project",
			UsageText: "project remove <name>",
			ArgsUsage: "name",
			Action:    remove,
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list projects",
			Action:  list,
		},
		{
			Name:      "open",
			Aliases:   []string{"o"},
			Usage:     "open project",
			UsageText: "project open <name>",
			ArgsUsage: "name",
			Action:    open,
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
