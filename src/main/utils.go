package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/montanaflynn/stats"
	"io"
	"log"
	"os"
	"strconv"
	"time"
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
	row := map[string]string{}
	var header []string
	var prevFdid int64 = 0
	if mysqlRows != nil {
		for mysqlRows.Next() {
			var id int64
			var fdid int64
			var fid int64
			var demographic_no int64
			var var_name string
			var var_value string
			mysqlRows.Scan(&id, &fdid, &fid, &demographic_no, &var_name, &var_value)
			if !isMember(var_name, header) {
				header = append(header, var_name)
			}
			row[var_name] = var_value
			if prevFdid == 0 {
				prevFdid = fdid
			}
			if fdid != prevFdid {
				rows = append(rows, row)
				prevFdid = fdid
			}
		}
	} else {
		os.Exit(2)
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
		fmt.Fprintln(v, *message)
		fmt.Fprintln(v, " ")
		varType := "string"
		counter := make(map[string]int)
		varNum := []float64{}
		for _, record := range csvMapValid {
			if n, err := strconv.ParseFloat(record[*message], 64); err == nil {
				varNum = append(varNum, n)
				varType = "num"
			} else {
				counter[record[*message]]++

			}
			// https://stackoverflow.com/questions/44417913/go-count-distinct-values-in-array-performance-tips
		}
		distinctStrings := make([]string, len(counter))
		i := 0
		for k := range counter {
			distinctStrings[i] = k
			i++
		}
		for _, s := range distinctStrings {
			fmt.Fprintln(v, s, " --> ", counter[s], " | ", counter[s]*100/recordCount, "%")
		}
		if varType == "num" {
			a, _ := stats.Sum(varNum)
			fmt.Fprintln(v, "Sum -->", a)
			a, _ = stats.Min(varNum)
			fmt.Fprintln(v, "Min -->", a)
			a, _ = stats.Max(varNum)
			fmt.Fprintln(v, "Max -->", a)
			a, _ = stats.Mean(varNum)
			fmt.Fprintln(v, "Mean -->", a)
			a, _ = stats.Median(varNum)
			fmt.Fprintln(v, "Median -->", a)
			a, _ = stats.StandardDeviation(varNum)
			fmt.Fprintln(v, "StdDev -->", a)

		}
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

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

func findDuplicates(csvMap []map[string]string) {
	var latest bool
	var included bool
	var demographicNo []string
	for _, v := range csvMap {
		latest = false
		included = true
		for k2, v2 := range v {
			if k2 == "eft_latest" && v2 == "1" {
				latest = true
			}
			if k2 == "dateCreated" {
				dateCreated, _ := time.Parse("2006-01-02", v2)
				_dateFrom, _ := time.Parse("2006-01-02", *dateFrom)
				_dateTo, _ := time.Parse("2006-01-02", *dateTo)
				if len(*dateFrom) > 0 && len(*dateTo) > 0 && !inTimeSpan(_dateFrom, _dateTo, dateCreated) {
					included = false
				}
			}
			if k2 == "demographic_no" {
				if !isMember(v2, demographicNo){
					demographicNo = append(demographicNo, v2)
					latest = true
				}
			}
			if *includeAll {
				latest = true
			}
		}
		if latest && !included {
			csvMapValid = append(csvMapValid, v)
			recordCount++
		}
	}
}