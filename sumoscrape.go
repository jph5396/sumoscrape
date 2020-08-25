package main

import (
	"fmt"
	"os"

	"github.com/jph5396/sumoscrape/commands"
	"github.com/jph5396/sumoscrape/sumoutils"
)

func main() {
	sumoutils.PrintTitle()

	cmdRegistry := []commands.Command{
		commands.NewBanzukeCommand(),
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

	//TODO: complete torikumi command and delete old code.

	// torikumiCommand := flag.NewFlagSet("torikumi", flag.ExitOnError)
	// torikumiBashoIdFlag := torikumiCommand.Int("basho-id", -1, "The basho to target. <YYYYMM>")
	// torikumiDayFlag := torikumiCommand.Int("day", -1, "day to target.")

	// switch os.Args[1] {
	// case "banzuke":
	// 	banzukeCommand.Parse(os.Args[2:])
	// 	list := commands.BanzukeScrape(*banzukeBashoIdFlag)
	// 	for _, item := range list {
	// 		item.PrintData()
	// 	}

	// case "torikumi":
	// 	torikumiCommand.Parse(os.Args[2:])
	// 	fmt.Printf("basho-id: %v. day: %v", *torikumiBashoIdFlag, *torikumiDayFlag)
	// 	fmt.Println()
	// default:
	// 	fmt.Println("Not a valid command.")
	// 	os.Exit(1)
	// }

}

//https://blog.rapid7.com/2016/08/04/build-a-simple-cli-tool-with-golang/
