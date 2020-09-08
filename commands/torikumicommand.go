package commands

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/gocolly/colly/v2"
	"github.com/jph5396/sumoscrape/sumomodel"
	"github.com/jph5396/sumoscrape/sumoutils"
)

type (

	// TorikumiCommand command and flagset to be executed.
	TorikumiCommand struct {
		TorikumiFlagSet *flag.FlagSet
		bashoID         int
		day             int
		divisions       DivisionFlag
		sysConfig       sumoutils.Config
	}
)

//NewTorikumiCommand returns a new ToeikumiCommand type
func NewTorikumiCommand(config sumoutils.Config) *TorikumiCommand {

	cmd := &TorikumiCommand{
		TorikumiFlagSet: flag.NewFlagSet("torikumi", flag.ExitOnError),
		sysConfig:       config,
	}

	cmd.TorikumiFlagSet.IntVar(&cmd.bashoID, "basho-id", -1, "The basho to target <YYYYMM>")
	cmd.TorikumiFlagSet.IntVar(&cmd.day, "day", -1, "the day to get bouts for must be a value between 1-16")
	cmd.TorikumiFlagSet.Var(&cmd.divisions, "division", "A division to target. Repeatable")

	return cmd
}

//CommandName returns command name.
func (cmd *TorikumiCommand) CommandName() string {
	return cmd.TorikumiFlagSet.Name()
}

//Parse parses command arguments and returns an error if any of the values are invalid.
func (cmd *TorikumiCommand) Parse(osArgs []string) error {
	cmd.TorikumiFlagSet.Parse(osArgs)
	if len(cmd.divisions) < 1 {
		cmd.divisions = append(cmd.divisions, "M", "J")
	}

	return nil
}

//Run logic to be executed when the torikumi command is called.
func (cmd *TorikumiCommand) Run() error {

	c := colly.NewCollector()
	var BoutList []sumomodel.Bout
	RequestedDivisions, err := sumomodel.GetDivisionList(cmd.divisions)
	if err != nil {
		fmt.Println(err.Error())
	}

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting", r.URL)
	})

	// check if website returned proper response. If it did not, inform and exit.
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("%v returned an HTTP code (%v) indicating the request failed. try again later.", r.Request.URL, r.StatusCode)
		fmt.Println()
		os.Exit(1)
	})

	c.OnHTML("table.tk_table", func(e *colly.HTMLElement) {

		e.ForEach("tr", func(i int, tr *colly.HTMLElement) {
			var NewBout sumomodel.Bout
			// set boutnum to the current iteration
			NewBout.Boutnum = i

			tr.ForEach("td", func(j int, td *colly.HTMLElement) {
				// we loop through the td tags on the table and and set the appropriate
				// variable based on the column index

				// the td at 0 represents the day
				if j == 0 {
					day, err := strconv.Atoi(td.Text)
					if err != nil {
						panic(err)
					}

					NewBout.Day = day
				}

				// td at 2 represents the division
				if j == 2 {
					rawdiv, err := strconv.Atoi(td.Text)
					if err != nil {
						panic(err)
					}

					// sites division id starts at 5 while ours start at 1
					NewBout.Division = rawdiv - 4
				}

				// td at 5 represents the first (east) wrestler
				if j == 5 {
					var shikonaTag sumomodel.ShikonaATag
					shikonaTag.ParseShikonaATag(td)
					NewBout.EastRikishiID = shikonaTag.Id
					NewBout.EastRikishiName = shikonaTag.Name
				}

				// td at 7 represents the EastWin variable
				if j == 7 {
					NewBout.EastWin = didWin(td.Text)
				}

				// td at 8 represents the kimarite
				if j == 8 {
					NewBout.Kimarite = td.Text
				}

				//td at 9 represents WestWin
				if j == 9 {
					NewBout.WestWin = didWin(td.Text)
				}

				// td at 11 represents west rikishi
				if j == 11 {
					var shikonaTag sumomodel.ShikonaATag
					shikonaTag.ParseShikonaATag(td)
					NewBout.WestRikishiID = shikonaTag.Id
					NewBout.WestRikishiName = shikonaTag.Name
				}
			})

			if IsRequestedDivisionByID(RequestedDivisions, NewBout.Division) {
				BoutList = append(BoutList, NewBout)
			}
		})

	})

	c.OnScraped(func(r *colly.Response) {
		fileName := sumoutils.CreateFileName(cmd.CommandName() + fmt.Sprintf("day%v", cmd.day))
		fmt.Println(cmd.sysConfig.SavePath)
		err := sumoutils.JSONFileWriter(cmd.sysConfig.SavePath+fileName, BoutList)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	})

	c.Visit(fmt.Sprintf("http://sumodb.sumogames.de/Results.aspx?b=%v&d=%v&simple=on", cmd.bashoID, cmd.day))

	return nil
}

func didWin(outcome string) bool {
	if outcome == "W" || outcome == "FS" {
		return true
	}
	return false
}
