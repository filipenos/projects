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
			Name:      "add",
			Aliases:   []string{"a"},
			Usage:     "add new project",
			UsageText: "project add <name> <path>",
			ArgsUsage: "name path",
			Action:    add,
		},
		{
			Name:   "current",
			Usage:  "add the current path as project",
			Action: addCurrent,
		},
		{
			Name:      "remove",
			Usage:     "remove project",
			UsageText: "project remove <name>",
			ArgsUsage: "name",
			Action:    remove,
		},
		{
			Name:   "list",
			Usage:  "list projects",
			Action: list,
		},
		{
			Name:      "open",
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

func log(msg string, args ...interface{}) {
	fmt.Printf("[projects] %s\n", fmt.Sprintf(msg, args...))
}
