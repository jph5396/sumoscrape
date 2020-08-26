package sumomodel

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

type (
	// Rikishi data struct.
	Rikishi struct {
		Id         int
		Name       string
		Rank       string
		Heya       string
		Shusshin   string
		HW         string
		Result     string
		Kanji      string
		Dob        string
		Firstbasho string
		Lastbasho  string
		Division   int
	}

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

// PrintData prints some of the rikishi structs data as a test.
func (r *Rikishi) PrintData() {
	fmt.Printf(
		"id: %v, rank: %v, name: %v, kanji: %v, heya: %v, shusshin: %v, dob = %v, results = %v",
		r.Id,
		r.Rank,
		r.Name,
		r.Kanji,
		r.Heya,
		r.Shusshin,
		r.Dob,
		r.Result)
	fmt.Println()
}

//ApplyTagResults apply tag results to Rikishi
func (r *Rikishi) ApplyTagResults(results ShikonaATag) {
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

	fmt.Println("received", element)
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
