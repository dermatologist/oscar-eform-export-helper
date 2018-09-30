package main

import (
	"github.com/jroimartin/gocui"
	"log"
	"os"
	"flag"
)

//This is how you declare a global variable
var csvMap []map[string] string

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true

	g.SetManagerFunc(layout)

	filePtr := flag.String("word", "test.csv", "The csv file to process")
	r, err := os.Open(*filePtr)
	if err != nil {
		log.Panicln(err)
	}
	csvMap = CSVToMap(r)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
