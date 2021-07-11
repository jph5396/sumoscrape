package commands

import (
	"fmt"
	"os"

	"github.com/jph5396/sumoscrape/sumoutils"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "Sumoscrape",
	Short: "gather sumo data from the web",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		sumoutils.PrintTitle()
		// note: all children inherit this command to make sure that a save dir is properly created.

		if cmd.Name() != "help" {
			directory := cmd.Flag("saveDir").Value.String()
			err := os.MkdirAll(directory, 0750)
			if err != nil {
				fmt.Printf("failed to create directory with error: %v", err.Error())
				os.Exit(1)
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
}

var SaveDir string = "temp/"

func Init() {

	RootCmd.PersistentFlags().StringVarP(&SaveDir, "saveDir", "s", "temp/", "Directory to save files to")
	RootCmd.PersistentFlags().StringVarP(&sumoutils.FileType, "FileType", "F", "json", "the file type the data should be saved in")
	RootCmd.AddCommand(
		NewBanzukeCommand(),
		NewTorikumiCommand())
}

func Execute(args []string) error {
	Init()
	err := RootCmd.Execute()
	return err
}
