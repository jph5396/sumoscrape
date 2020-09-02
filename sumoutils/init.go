package sumoutils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mitchellh/go-homedir"
)

type (

	//Config struct is used to store config vars. currently only contains 1 var
	// but declaring it for future use.
	Config struct {
		// the Path that data should be saved to.
		// will default to the config's directory.
		SavePath string `json:"SavePath"`
		Test     string `json:"Test"`
	}
)

// Init loads the config file into the system. Will prompt the user
// to create a new one if it does not exist.
func Init() Config {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println("home dir error.")
		os.Exit(1)
	}

	// Checks if config exists. If it doesnt it will create it.
	findOrCreateConfig(home)
	var config Config

	configFile, err := os.Open(home + "/sumoscrape/config.json")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer configFile.Close()
	configBytes, readerErr := ioutil.ReadAll(configFile)
	if readerErr != nil {
		fmt.Println(readerErr.Error())
	}

	marsh := json.Unmarshal(configBytes, &config)
	if marsh != nil {
		fmt.Println(marsh.Error())
	}

	return config

}

// returns config file path. If it does not exist, it will create and
// recursively return file path.
func findOrCreateConfig(home string) {

	configFilePath := fmt.Sprintf("%v/%v", home, "sumoscrape/config.json")
	_, err := os.Stat(configFilePath)
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

		var config Config

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

		config.Test = "test"
		//save config file to user home/sumoscrape/config.json
		jsonErr := JSONFileWriter(home+"/sumoscrape/config.json", config)
		if jsonErr != nil {
			fmt.Println(err.Error())
		}

		// recursively call findOrCreateConfig to return the file info.
		findOrCreateConfig(home)

		// if there is an error but it is not caused by a file not existing, notify and exit.
	} else if err != nil {
		fmt.Println("An unknown error occured. Cannot continue.")
		fmt.Printf("error: %v", err.Error())
		fmt.Println("")
	}
}
