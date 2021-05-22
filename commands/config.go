package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/jph5396/sumoscrape/sumoutils"
	"github.com/mitchellh/go-homedir"
)

type ConfigCommand struct {
	sysConfig sumoutils.Config
}

func NewConfigCommand(sysconf sumoutils.Config) *ConfigCommand {
	return &ConfigCommand{
		sysConfig: sysconf,
	}
}

func (cmd *ConfigCommand) CommandName() string {
	return "config"
}

func (cmd *ConfigCommand) Parse(args []string) error {
	if len(args) != 0 {
		return errors.New("config should not receive any args")
	}

	return nil
}

func (cmd *ConfigCommand) Run() error {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println("home dir error.")
		os.Exit(1)
	}

	configFilePath := fmt.Sprintf("%v/%v", home, "sumoscrape/config.json")
	info, err := os.Stat(configFilePath)
	if info != nil {
		fmt.Printf("Config file found at %v . do you want to overwrite? [Y for yes]", info.Name())
		var res string
		fmt.Scanln(&res)
		if res == "Y" {

			os.Remove(configFilePath)

			var config sumoutils.Config

			// prompt user for output directory. default to user home path /sumoscrape
			fmt.Println(fmt.Sprintf("Enter path to directory that should be saved to [%v]:", home+"/sumoscrape/data"))
			fmt.Scanln(&config.SavePath)

			//save config file to user home/sumoscrape/config.json
			jsonErr := sumoutils.JSONFileWriter(configFilePath, config)
			if jsonErr != nil {
				fmt.Println(err.Error())
			}

		} else {
			fmt.Println("Y not recevied. Cancelling ")
			return nil
		}
	}
	if os.IsNotExist(err) {
		fmt.Println("Config file not found... creating")

		// create sumoscrape directory with user only write access.
		err := os.Mkdir(home+"/sumoscrape", 0755)

		// Do not continue if an error occurs and it is not because the directory already exists.
		if err != nil && !os.IsExist(err) {

			fmt.Println("There was an error when creating the sumoscrape directory. Can not continue.")
			fmt.Printf("error: %v", err.Error())
			os.Exit(1)
		}

		var config sumoutils.Config

		// prompt user for output directory. default to user home path /sumoscrape
		fmt.Println(fmt.Sprintf("Enter path to directory that should be saved to [%v]:", home+"/sumoscrape/data"))
		fmt.Scanln(&config.SavePath)

		// we call mkdirall here because it does nothing if the directory already exists.
		if config.SavePath == "" {
			err := os.MkdirAll(home+"/sumoscrape/data", 0755)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			config.SavePath = home + "/sumoscrape/data/"
		} else {
			err := os.MkdirAll(config.SavePath, 0755)
			if err != nil {
				fmt.Printf("%v could not be created.", config.SavePath)
				fmt.Println(" " + err.Error())
			}
		}

		//save config file to user home/sumoscrape/config.json
		jsonErr := sumoutils.JSONFileWriter(home+"/sumoscrape/config.json", config)
		if jsonErr != nil {
			fmt.Println(err.Error())
		}

		// if there is an error but it is not caused by a file not existing, notify and exit.
	} else if err != nil {
		fmt.Println("An unknown error occured. Cannot continue.")
		fmt.Printf("error: %v", err.Error())
		fmt.Println("")
	}
	return nil
}
