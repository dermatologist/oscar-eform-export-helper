package main

import (
	"encoding/csv"
	"io"
	"log"
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

func mainOutput(v *gocui.View, message *string)  {
	fmt.Fprintf(v, "%s", *message)
	v.Editable = true
	v.Wrap = true
}

func sideOutput(v *gocui.View)  {
	//fmt.Fprintln(v, "Column 1")
	//fmt.Fprintln(v, "Column 2")
	//fmt.Fprintln(v, "Column 3")
	//fmt.Fprint(v, "\rWill be")
	//fmt.Fprint(v, "deleted\rColumn 4\nColumn 5")
	firstRecord := csvMap[0]
	for key, _ := range firstRecord {
		fmt.Fprintln(v, key)
	}
}