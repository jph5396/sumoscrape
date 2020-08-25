package sumomodel

import "fmt"

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
