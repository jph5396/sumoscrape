package commands

import "strings"

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
func (d *DivisionFlag) Set(str string) error {
	*d = append(*d, str)
	return nil
}
