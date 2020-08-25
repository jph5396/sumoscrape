package commands

//TODO: way too much sumo logic is stored in this command.
// this should be refactored to seperate that logic out into
// its own package.
import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/jph5396/sumoscrape/sumomodel"
)

type (
	//BanzukeCommand struct containing the Flags for the command and variables that the are used when parsing.
	BanzukeCommand struct {
		BanzukeFlags *flag.FlagSet
		// ID of target basho. in YYYYMM format
		bashoID int
	}
)

// NewBanzukeCommand creates Banzuke Command and flagset.
func NewBanzukeCommand() *BanzukeCommand {
	cmd := &BanzukeCommand{
		BanzukeFlags: flag.NewFlagSet("banzuke", flag.ExitOnError),
	}
	cmd.BanzukeFlags.IntVar(&cmd.bashoID, "basho-id", -1, "The basho to target <YYYYMM>")
	return cmd
}

// CommandName returns the name of the command
func (cmd *BanzukeCommand) CommandName() string {
	return cmd.BanzukeFlags.Name()
}

// Parse the args received from the OS
func (cmd *BanzukeCommand) Parse(osArgs []string) error {
	cmd.BanzukeFlags.Parse(osArgs)
	return nil
}

// Run runs the BanzukeCommand by reaching out to the target URL and parsing the tables representing the banzuke.
func (cmd *BanzukeCommand) Run() error {
	if cmd.bashoID == -1 {
		fmt.Println("the --basho-id flag must be set in YYYYMM format")
		os.Exit(1)
	}

	c := colly.NewCollector()
	RikishiList := []sumomodel.Rikishi{}

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting", r.URL)
	})

	// check if website returned proper response. If it did not, inform and exit.
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("%v returned an HTTP code (%v) indicating the request failed. try again later.", r.Request.URL, r.StatusCode)
		fmt.Println()
		os.Exit(1)
	})

	c.OnHTML("table.banzuke", func(e *colly.HTMLElement) {

		tableCaption := e.ChildText("caption")
		// only target Makuuchi and juryo divisions.
		if strings.Contains(tableCaption, "Makuuchi") || strings.Contains(tableCaption, "Juryo") {

			// each tr should represent 1 rikishi
			e.ForEach("tr", func(i int, tr *colly.HTMLElement) {

				var newRikishi sumomodel.Rikishi

				//using td index to identify column.
				tr.ForEach("td", func(j int, td *colly.HTMLElement) {
					if j == 0 {
						newRikishi.Rank = td.Text

						//set division based on rank.
						if strings.Contains(td.Text, "J") {
							newRikishi.Division = 2
						} else {
							newRikishi.Division = 1
						}
					}
					if j == 1 {
						aTagResults := parseShikonaATag(td)
						newRikishi.ApplyTagResults(aTagResults)
					}
					if j == 2 {
						newRikishi.Result = td.Text
					}
				})

				if newRikishi.Id != 0 {
					RikishiList = append(RikishiList, newRikishi)
				}

			})
		}

	})

	c.Visit(fmt.Sprintf("http://sumodb.sumogames.de/Banzuke.aspx?b=%v&hl=on&c=on", cmd.bashoID))

	for _, rikishi := range RikishiList {
		rikishi.PrintData()
	}

	return nil
}

// function should parse the title and href from the A tag and return a ShikonaATag struct
func parseShikonaATag(element *colly.HTMLElement) sumomodel.ShikonaATag {
	var returnVal sumomodel.ShikonaATag

	titleArr := strings.Split(element.ChildAttr("a", "title"), ",")
	newid, err := strconv.Atoi(strings.Split(element.ChildAttr("a", "href"), "=")[1])
	if err != nil {
		panic(err)
	}

	returnVal.Id = newid
	returnVal.Name = element.Text
	returnVal.Kanji = titleArr[0]
	returnVal.Heya = titleArr[1]
	returnVal.Shusshin = titleArr[2]
	returnVal.Dob = titleArr[3]
	returnVal.Firstbasho = titleArr[4]
	returnVal.Lastbasho = titleArr[5]
	returnVal.HW = titleArr[6]

	return returnVal
}
