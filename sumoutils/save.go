package sumoutils

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jszwec/csvutil"
)

var FileType string = "json"

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

func CSVFileWriter(path string, data interface{}) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	encoder := csvutil.NewEncoder(writer)
	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	writer.Flush()
	return writer.Error()
}

// CreateFileName creates a file name with the supplied prefix followed by the
// time in RFC822 format
func CreateFileName(prefix string) string {
	return strings.ReplaceAll(prefix+time.Now().Format("02 Jan 06"), " ", "_") + "." + FileType
}

// SaveFile takes in the file path and data to save to the file. selects the appropriate
// file type based on the fileType variable.
func SaveFile(path string, data interface{}) error {
	var err error = nil
	switch FileType {
	case "csv":
		err = CSVFileWriter(path, data)
	default:
		err = JSONFileWriter(path, data)
	}
	return err
}
