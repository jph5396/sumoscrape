package commands

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/jph5396/sumomodel"
	"github.com/jph5396/sumoscrape/sumoutils"
	"github.com/spf13/cobra"
)

var dayID int

var torikumiCommand = &cobra.Command{
	Use:   "torikumi",
	Short: "The torikumi to get data from",
	Long:  "The torikumi command scrapes all bouts for a given day.",
	PreRun: func(cmd *cobra.Command, args []string) {
		if strings.Contains(cmd.Flag("division").Value.String(), "all") {
			divisions = []string{"M", "J", "Ms", "Sd", "Jd", "Jk"}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		c := colly.NewCollector()
		var BoutList []sumomodel.Bout
		RequestedDivisions, err := sumomodel.GetDivisionList(divisions)
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
				NewBout.BashoID = bashoID

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
						var shikonaTag ShikonaATag
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
						var shikonaTag ShikonaATag
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
			fileName := sumoutils.CreateFileName(cmd.Name())
			dir := cmd.Flag("saveDir").Value.String()
			if string(dir[len(dir)-1:]) != "/" {
				dir = dir + "/"
			}

			err := sumoutils.SaveFile(dir+fileName, BoutList)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		})

		c.Visit(fmt.Sprintf("http://sumodb.sumogames.de/Results.aspx?b=%v&d=%v&simple=on", bashoID, dayID))
	},
}

func NewTorikumiCommand() *cobra.Command {
	torikumiCommand.Flags().IntVarP(&bashoID, "basho-id", "b", 0, "The Basho to get data for YYYYMM")
	torikumiCommand.Flags().IntVar(&dayID, "day", 0, "the day to get data from")
	torikumiCommand.Flags().StringArrayVarP(&divisions, "division", "d", []string{"M", "J"}, "The Divisions to target. options: M, J ,Ms, Sd, Jd, Jk, or all")
	torikumiCommand.MarkFlagRequired("basho-id")
	torikumiCommand.MarkFlagRequired("day")

	return torikumiCommand
}

func didWin(outcome string) bool {
	if outcome == "W" || outcome == "FS" {
		return true
	}
	return false
}
