package commands

import (
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/jph5396/sumomodel"
)

// common functions and types that are used by multiple commands.

type (
	//DivisionFlag an array of strings that are used to decide which divisions to target when scrapping.
	// it implements the value interface from the flag package.
	DivisionFlag []string

	// ShikonaATag the anchor elements store a portion of desired data
	// in the title. This struct allows that data to be easily passed out
	// of function and added to rikishi.
	ShikonaATag struct {
		Id         int
		Name       string
		Kanji      string
		Heya       string
		Shusshin   string
		Dob        string
		Firstbasho string
		Lastbasho  string
		HW         string
	}
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
// by the user. it will return the division ID and true if it was requested.
func IsRequestedDivision(d []sumomodel.Division, str string) (int, bool) {
	for _, item := range d {
		if str == item.DivLongName {
			return item.ID, true
		}
	}
	return -1, false
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

//ApplyTagResults apply tag results to Rikishi
func ApplyTagResults(results ShikonaATag, r *sumomodel.Rikishi) {
	r.Id = results.Id
	r.Name = results.Name
	r.Kanji = results.Kanji
	r.Heya = results.Heya
	r.Shusshin = results.Shusshin
	r.Dob = results.Dob
	r.HW = results.HW
	r.Firstbasho = results.Firstbasho
	r.Lastbasho = results.Lastbasho
}

//ParseShikonaATag takes a colly HtmlElement (which should be a td that contains an a tag with a title)
// then parses its contents and applies them to a shikonaATag
func (s *ShikonaATag) ParseShikonaATag(element *colly.HTMLElement) {

	titleArr := strings.Split(element.ChildAttr("a", "title"), ",")
	newid, err := strconv.Atoi(strings.Split(element.ChildAttr("a", "href"), "=")[1])
	if err != nil {
		panic(err)
	}

	s.Id = newid
	s.Name = element.Text
	s.Kanji = titleArr[0]
	s.Heya = titleArr[1]
	s.Shusshin = titleArr[2]
	s.Dob = titleArr[3]
	s.Firstbasho = titleArr[4]
	s.Lastbasho = titleArr[5]
	s.HW = titleArr[6]
}
