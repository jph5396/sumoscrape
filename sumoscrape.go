package main

import (
	"fmt"
	"os"

	"github.com/jph5396/sumoscrape/commands"
	"github.com/jph5396/sumoscrape/sumoutils"
)

func main() {
	sumoutils.PrintTitle()

	config := sumoutils.Init()

	cmdRegistry := []commands.Command{
		commands.NewBanzukeCommand(config),
		commands.NewTorikumiCommand(config),
		commands.NewConfigCommand(config),
	}

	if len(os.Args) < 2 {
		fmt.Println("a subcommand is required. the available subcommands are banzuke and torikumi.")
		os.Exit(1)
	}

	// TODO: add error handling.
	for _, cmd := range cmdRegistry {
		if cmd.CommandName() == os.Args[1] {
			cmd.Parse(os.Args[2:])
			cmd.Run()
		}
	}
}

//https://blog.rapid7.com/2016/08/04/build-a-simple-cli-tool-with-golang/
//https://www.digitalocean.com/community/tutorials/how-to-use-the-flag-package-in-go
//https://golang.org/pkg/flag/
