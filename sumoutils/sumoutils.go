// Package sumoutils is a collection of utility functions written for sumo-scrape
// not all are super important. Some are just for fun!
package sumoutils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/common-nighthawk/go-figure"
)

// PrintTitle Prints ANSCII art title when program is ran.
func PrintTitle() {
	rand.Seed(time.Now().UnixNano())
	var colorList = []string{"red", "purple", "cyan", "blue", "green"}
	var fontList = []string{"smkeyboard", "smshadow", "small", "wavy", "weird"}

	titleFig := figure.NewColorFigure(
		"Sumo-Scrape",
		fontList[rand.Intn(len(fontList))],
		colorList[rand.Intn(len(colorList))],
		true)
	titleFig.Print()
}

// JSONFileWriter reusable function for writing structs to json files.
// will return any errors that occur, or nil if it succeeds.
func JSONFileWriter(path string, data interface{}) error {

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	encodeErr := encoder.Encode(data)
	if encodeErr != nil {
		return encodeErr
	}

	defer file.Close()
	fmt.Println("Created file: ", path)
	return nil
}

// CreateFileName creates a file name in the <command name>YYYY-MM-DD
func CreateFileName(cmd string) string {
	return strings.ReplaceAll(cmd+time.Now().String(), " ", "") + ".json"
}
