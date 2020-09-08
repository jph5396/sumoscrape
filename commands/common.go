package commands

import (
	"strings"

	"github.com/jph5396/sumoscrape/sumomodel"
)

// common functions and types that are used by multiple commands.

type (
	//DivisionFlag an array of strings that are used to decide which divisions to target when scrapping.
	// it implements the value interface from the flag package.
	DivisionFlag []string
)

func (d *DivisionFlag) String() string {
	var builder strings.Builder

	for _, item := range *d {
		builder.WriteString(item)
	}
	return builder.String()
}

// Set Division flag implementation of the Set(string)
// function required by the flags.Value interface
func (d *DivisionFlag) Set(s string) error {
	*d = append(*d, s)
	return nil
}

// IsRequestedDivision is used to check if a banzuke/ bout was requested
// by the user. it will return true if the str is present in the division
// list.
func IsRequestedDivision(d []sumomodel.Division, str string) bool {
	for _, item := range d {
		if str == item.DivLongName {
			return true
		}
	}
	return false
}

// IsRequestedDivisionByID is used to check if a banzuke/ bout was requested
// by the user based on the ID provided.
func IsRequestedDivisionByID(d []sumomodel.Division, id int) bool {
	for _, item := range d {
		if id == item.ID {
			return true
		}
	}
	return false
}
