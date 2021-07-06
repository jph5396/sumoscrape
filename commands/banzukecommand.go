package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/jph5396/sumomodel"
	"github.com/jph5396/sumoscrape/sumoutils"
	"github.com/spf13/cobra"
)

var bashoID int
var divisions []string

var banzukeCommand = &cobra.Command{
	Use:   "banzuke",
	Short: "get Banzuke data",
	Long:  "Get Banzuke Data from sumodb",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if bashoID == 0 {
			fmt.Println("the --basho-id flag must be set in YYYYMM format")
			os.Exit(1)
		}
		dir := cmd.Flag("saveDir").Value.String()
		c := colly.NewCollector()
		RikishiList := []sumomodel.Rikishi{}
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

		c.OnHTML("table.banzuke", func(e *colly.HTMLElement) {

			tableDivision := strings.Split(e.ChildText("caption"), " ")[0]

			// check if the current division was requested.
			if divisionID, ok := IsRequestedDivision(RequestedDivisions, tableDivision); ok {

				// each tr should represent 1 rikishi
				e.ForEach("tr", func(i int, tr *colly.HTMLElement) {

					var newRikishi sumomodel.Rikishi
					newRikishi.BashoID = bashoID
					newRikishi.Division = divisionID

					//using td index to identify column.
					tr.ForEach("td", func(j int, td *colly.HTMLElement) {
						if j == 0 {
							newRikishi.Rank = td.Text
						}
						if j == 1 {
							var aTag ShikonaATag
							aTag.ParseShikonaATag(td)
							ApplyTagResults(aTag, &newRikishi)
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

		// save data post scrape.
		c.OnScraped(func(r *colly.Response) {
			fileName := sumoutils.CreateFileName(cmd.Name())

			if string(dir[len(dir)-1:]) != "/" {
				dir = dir + "/"
			}

			err := sumoutils.JSONFileWriter(dir+fileName, RikishiList)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

		})

		c.Visit(fmt.Sprintf("http://sumodb.sumogames.de/Banzuke.aspx?b=%v&hl=on&c=on", bashoID))

	},
}

func NewBanzukeCommand() *cobra.Command {
	banzukeCommand.Flags().IntVarP(&bashoID, "basho-id", "b", 0, "The Basho ID to get the data for YYYYMM")
	banzukeCommand.Flags().StringArrayVarP(&divisions, "division", "d", []string{"M", "J"}, "The Division to target. Repeatable <M,J,Ms,Sd, Jd, Jk>")
	banzukeCommand.MarkFlagRequired("basho-id")
	return banzukeCommand
}
