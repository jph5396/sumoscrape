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

	var config Config

	configFile, err := os.Open(home + "/sumoscrape/config.json")
	if err == nil {
		configBytes, readerErr := ioutil.ReadAll(configFile)
		if readerErr != nil {
			fmt.Println(readerErr.Error())
		}

		marsh := json.Unmarshal(configBytes, &config)
		if marsh != nil {
			fmt.Println(marsh.Error())
		}
	}
	if os.IsNotExist(err) {
		fmt.Println("No config file found. using current directory")
		fmt.Println("run sumoscrape config to create a config file")

		wdir, err := os.Getwd()
		if err != nil {
			fmt.Println(err.Error())
		}

		config.SavePath = wdir + "/temp/"
		os.Mkdir("temp", 0755)
	} else if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return config

}
