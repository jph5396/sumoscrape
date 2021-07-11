// Package sumoutils is a collection of utility functions written for sumo-scrape
// not all are super important. Some are just for fun!
package sumoutils

import (
	"math/rand"
	"time"

	"github.com/common-nighthawk/go-figure"
)

// PrintTitle Prints ANSCII art title when program is ran.
func PrintTitle() {
	rand.Seed(time.Now().UnixNano())
	var colorList = []string{"red", "purple", "cyan", "blue", "green"}
	var fontList = []string{"smkeyboard", "smshadow", "small", "wavy", "weird"}

	titleFig := figure.NewColorFigure(
		"Sumo-Scrape",
		fontList[rand.Intn(len(fontList))],
		colorList[rand.Intn(len(colorList))],
		true)
	titleFig.Print()
}
