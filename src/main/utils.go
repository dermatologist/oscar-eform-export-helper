package main

import (
	"encoding/csv"
	"io"
	"log"
	"flag"
	"io/ioutil"
	"fmt"
	"github.com/jroimartin/gocui"
)

// CSVToMap takes a reader and returns an array of dictionaries, using the header row as the keys
func CSVToMap(reader io.Reader) []map[string]string {
	r := csv.NewReader(reader)
	rows := []map[string]string{}
	var header []string
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if header == nil {
			header = record
		} else {
			dict := map[string]string{}
			for i := range header {
				dict[header[i]] = record[i]
			}
			rows = append(rows, dict)
		}
	}
	return rows
}

func mainOutput(v *gocui.View)  {
	filePtr := flag.String("word", "test.csv", "The csv file to process")
	b, err := ioutil.ReadFile(*filePtr)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(v, "%s", b)
	v.Editable = true
	v.Wrap = true
}