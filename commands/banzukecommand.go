package commands

//TODO: way too much sumo logic is stored in this command.
// this should be refactored to seperate that logic out into
// its own package.
import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

type (
	// Rikishi data struct.
	Rikishi struct {
		id         int
		name       string
		rank       string
		heya       string
		shusshin   string
		hw         string
		result     string
		kanji      string
		dob        string
		firstbasho string
		lastbasho  string
		division   int
	}

	// ShikonaATag the anchor elements store a portion of desired data
	// in the title. This struct allows that data to be easily passed out
	// of function and added to rikishi.
	ShikonaATag struct {
		id         int
		name       string
		kanji      string
		heya       string
		shusshin   string
		dob        string
		firstbasho string
		lastbasho  string
		hw         string
	}
)

// BanzukeScrape will scrape the banzuke of the given basho.
func BanzukeScrape(basho int) []Rikishi {
	if basho == -1 {
		fmt.Println("the --basho-id flag must be set in YYYYMM format")
		os.Exit(1)
	}

	c := colly.NewCollector()
	MasterBanzukeList := []Rikishi{}

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting", r.URL)
	})

	// check if website returned proper response. If it did not, inform and exit.
	c.OnResponse(func(r *colly.Response) {
		if r.StatusCode != 200 {
			fmt.Printf("%v returned an HTTP code (%v) indicating the request failed. try again later.", r.Request.URL, r.StatusCode)
			fmt.Println()
			os.Exit(1)
		}
	})

	c.OnHTML("table.banzuke", func(e *colly.HTMLElement) {

		tableCaption := e.ChildText("caption")
		// only target Makuuchi and juryo divisions.
		if strings.Contains(tableCaption, "Makuuchi") || strings.Contains(tableCaption, "Juryo") {

			// each tr should represent 1 or 2 rikishi.
			e.ForEach("tr", func(i int, tr *colly.HTMLElement) {

				var EastRikishi Rikishi
				var WestRikishi Rikishi

				tr.ForEach("td", func(j int, td *colly.HTMLElement) {
					if containsClass("shikona", td) || containsClass("debut", td) {
						namedata := parseShikonaATag(td)
						if identifySide(j) == "east" {
							EastRikishi.applyTagResults(namedata)
						} else {
							WestRikishi.applyTagResults(namedata)
						}
					}
					if containsClass("short_rank", td) {
						var baseRank string
						if isSanyaku(td.Text) {
							baseRank = getSanyakuRank(td.Text, MasterBanzukeList)
						} else {
							baseRank = td.Text
						}

						EastRikishi.rank = baseRank + "e"
						WestRikishi.rank = baseRank + "w"

						// set division based on rank.
						if strings.Contains(td.Text, "J") {
							EastRikishi.division = 2
							WestRikishi.division = 2
						} else {
							EastRikishi.division = 1
							WestRikishi.division = 1
						}
					}
					if containsClass("", td) {
						side := identifySide(j)
						if side == "east" {
							EastRikishi.result = td.Text
						} else {
							WestRikishi.result = td.Text
						}
					}

					// empty cell indicates there is not a rikishi for that side.
					// identify that side and set the corresponding rikishi id to zeros
					if containsClass("emptycell", td) {
						side := identifySide(j)
						if side == "east" {
							EastRikishi.id = 0
						} else {
							WestRikishi.id = 0
						}
					}

				})
				if EastRikishi.id != 0 {
					MasterBanzukeList = append(MasterBanzukeList, EastRikishi)
				}
				if WestRikishi.id != 0 {
					MasterBanzukeList = append(MasterBanzukeList, WestRikishi)
				}

			})
		}

	})

	c.Visit(fmt.Sprintf("http://sumodb.sumogames.de/Banzuke.aspx?b=%v", basho))

	return MasterBanzukeList
}

func identifySide(i int) string {
	if i == 1 || i == 0 {
		return "east"
	}
	return "west"
}

// check if rank is sanyaku
func isSanyaku(rank string) bool {
	sanyakuRanks := []string{"Y", "O", "S", "K"}
	for _, item := range sanyakuRanks {
		if item == rank {
			return true
		}
	}
	return false
}

// because sanyaku ranks to do not have their number, we need to add it.
// we check the master list to see what the previous rank was set to.
func getSanyakuRank(rank string, masterList []Rikishi) string {

	// if list has no members in it, return a "{rank}1"
	if len(masterList) == 0 {
		return rank + "1"

		// if the last rikishi rank is equal to rank, get last rikishi number and add one for new rank.
	} else if lastRikishi := masterList[len(masterList)-1]; string(lastRikishi.rank[0]) == rank {
		regex := regexp.MustCompile("[0-9]+")
		oldRankNum, err := strconv.Atoi(regex.FindString(lastRikishi.rank))
		if err != nil {
			panic(err)
		}

		return rank + strconv.Itoa(oldRankNum+1)
	}

	// rank is not equal.
	return rank + "1"
}

// uses querySelectory to determine if item exists.
func containsClass(class string, collyEl *colly.HTMLElement) bool {
	if collyEl.Attr("class") == class {
		return true
	}
	return false
}

// function should parse the title and href from the A tag and return a ShikonaATag struct
func parseShikonaATag(element *colly.HTMLElement) ShikonaATag {
	var returnVal ShikonaATag

	titleArr := strings.Split(element.ChildAttr("a", "title"), ",")
	newid, err := strconv.Atoi(strings.Split(element.ChildAttr("a", "href"), "=")[1])
	if err != nil {
		panic(err)
	}

	returnVal.id = newid
	returnVal.name = element.Text
	returnVal.kanji = titleArr[0]
	returnVal.heya = titleArr[1]
	returnVal.shusshin = titleArr[2]
	returnVal.dob = titleArr[3]
	returnVal.firstbasho = titleArr[4]
	returnVal.lastbasho = titleArr[5]
	returnVal.hw = titleArr[6]

	return returnVal
}

// apply tag results to Rikishi
func (r *Rikishi) applyTagResults(results ShikonaATag) {
	r.id = results.id
	r.name = results.name
	r.kanji = results.kanji
	r.heya = results.heya
	r.shusshin = results.shusshin
	r.dob = results.dob
	r.hw = results.hw
	r.firstbasho = results.firstbasho
	r.lastbasho = results.lastbasho
}

// PrintData prints some of the rikishi structs data as a test.
func (r *Rikishi) PrintData() {
	fmt.Printf(
		"id: %v, rank: %v, name: %v, kanji: %v, heya: %v, shusshin: %v, dob = %v, results = %v",
		r.id,
		r.rank,
		r.name,
		r.kanji,
		r.heya,
		r.shusshin,
		r.dob,
		r.result)
	fmt.Println()
}
