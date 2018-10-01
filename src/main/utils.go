package main

import (
	"encoding/csv"
	"fmt"
	"github.com/jroimartin/gocui"
	"io"
	"log"
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

func mainOutput(g *gocui.Gui, message *string) {
	if v, err := g.SetCurrentView("main"); err != nil {
		log.Panicln(err)
	} else {
		v.Editable = true
		v.Wrap = true
		v.Clear()
		fmt.Fprintf(v, "%s", *message)
		g.SetCurrentView("side")
		recover()
	}
}

func sideOutput(g *gocui.Gui) {
	toIgnore := []string{"id", "fdid", "dateCreated", "eform_link", "StaffSig", "SubmitButton", "efmfid"}
	if v, err := g.SetCurrentView("side"); err != nil {
		log.Panicln(err)
	} else {
		firstRecord := csvMap[0]
		for key, _ := range firstRecord {
			if !isMember(key, toIgnore) {
				fmt.Fprintln(v, key)
			}
		}
	}
}

func isMember(s string, a []string) bool {
	for _, v := range a {
		if v == s {
			return true
		}
	}
	return false
}
