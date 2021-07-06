package main

import (
	"fmt"
	"os"

	"github.com/jph5396/sumoscrape/commands"
)

func main() {

	err := commands.Execute(os.Args[1:])
	if err != nil {
		fmt.Println(err.Error())
	}
}

//https://blog.rapid7.com/2016/08/04/build-a-simple-cli-tool-with-golang/
//https://www.digitalocean.com/community/tutorials/how-to-use-the-flag-package-in-go
//https://golang.org/pkg/flag/
