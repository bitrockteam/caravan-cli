package main

import (
	"caravan/commands"
	"log"
	"os"

	"github.com/mitchellh/cli"
)

func main() {

	wd := ".caravan"
	//TODO metadata handling
	repos := []string{"caravan", "caravan-baking", "caravan-infra", "caravan-platform", "caravan-application-support"}

	c := cli.NewCLI("caravancli", "0.0.1")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"init": func() (cli.Command, error) {
			return &commands.InitCommand{
					Workdir: wd,
					Repos:   repos,
				},
				nil
		},
		"clean": func() (cli.Command, error) {
			return &commands.CleanupCommand{
					Workdir: wd,
				},
				nil
		},
		"bake": func() (cli.Command, error) {
			return &commands.BakeVMCommand{
					Workdir: wd,
				},
				nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
