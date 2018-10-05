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

func MysqlToMap(mysqlRows *sql.Rows) []map[string]string {
	rows := []map[string]string{}
	var header []string
	for mysqlRows.Next() {
		var id int64
		var fdid int64
		var fid int64
		var demographic_no int64
		var var_name string
		var var_value string
		mysqlRows.Scan(&id, &fdid, &fid, &demographic_no, &var_name, &var_value)
		if !isMember(var_name, header){
			header = append(header, var_name)
		}
		rows[fdid][var_name] = var_value
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
