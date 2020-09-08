package sumomodel

import "errors"

type (
	//Division Represents the a division of sumo wrestlers.
	// Each division contains an ID, the full name, a short version
	// of the name, and a Regex string used to confirm if a rank is
	// in said division.
	Division struct {
		ID          int
		DivLongName string
		ShortName   string
		RankRegex   string
	}
)

// Divisions are not exported and can only be accessed via functions.
var makuuchi = Division{
	ID:          1,
	DivLongName: "Makuuchi",
	ShortName:   "M",
	RankRegex:   "[YOSKM]\\d{1,3}[ew](HD|YO)?",
}

var juryo = Division{
	ID:          2,
	DivLongName: "Juryo",
	ShortName:   "J",
	RankRegex:   "J\\d{1,3}[ew]",
}

var makushita = Division{
	ID:          3,
	DivLongName: "Makushita",
	ShortName:   "Ms",
	RankRegex:   "Ms\\d{1,3}([ew]|TD)?",
}

var sandanme = Division{
	ID:          4,
	DivLongName: "Sandanme",
	ShortName:   "Sd",
	RankRegex:   "Sd\\d{1,4}([ew]|TD)?",
}

var jonidan = Division{
	ID:          5,
	DivLongName: "Jonidan",
	ShortName:   "Jd",
	RankRegex:   "Jd\\d{1,4}[ew]",
}

var jonokuchi = Division{
	ID:          6,
	DivLongName: "Jonokuchi",
	ShortName:   "Jk",
	RankRegex:   "Jk\\d{1,3}[ew]",
}

// create maps based on both long and short forms of the divisions.
var shortForm = map[string]Division{
	makuuchi.ShortName:  makuuchi,
	juryo.ShortName:     juryo,
	makushita.ShortName: makushita,
	sandanme.ShortName:  sandanme,
	jonidan.ShortName:   jonidan,
	jonokuchi.ShortName: jonokuchi,
}
var longForm = map[string]Division{
	makuuchi.DivLongName:  makuuchi,
	juryo.DivLongName:     juryo,
	makushita.DivLongName: makushita,
	sandanme.DivLongName:  sandanme,
	jonidan.DivLongName:   jonidan,
	jonokuchi.DivLongName: jonokuchi,
}

// GetDivision returns the requested division if it exists, or nil if it doesnt.
// Divisions can be requested by either their short or long forms
func GetDivision(str string) (Division, error) {
	// check if str is less than 3 characters. If it is,
	// get Division via the short form.
	if len(str) < 3 {
		if val, ok := shortForm[str]; ok {
			return val, nil
		}
		return Division{}, errors.New("division does not exist")
	}

	if val, ok := longForm[str]; ok {
		return val, nil
	}

	return Division{}, errors.New("division does not exist")
}

//GetDivisionList takes in a list of requested divisions and
func GetDivisionList(list []string) ([]Division, error) {

	var divList []Division
	for _, item := range list {
		division, err := GetDivision(item)
		if err != nil {
			return nil, err
		}
		divList = append(divList, division)
	}

	return divList, nil
}
