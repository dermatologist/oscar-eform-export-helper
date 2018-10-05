package main

import (
	"flag"
	"github.com/jroimartin/gocui"
	"log"
	"os"
)

//This is how you declare a global variable
var csvMap []map[string]string
var csvMapValid []map[string]string
var recordCount int

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true

	g.SetManagerFunc(layout)

	// ./command-line-flags -word=opt -numb=7 -fork -svar=flag
	filePtr := flag.String("file", "test.csv", "The csv file to process")
	r, err := os.Open(*filePtr)
	if err != nil {
		log.Panicln(err)
	}
	csvMap = CSVToMap(r)
	findDuplicates(csvMap)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
