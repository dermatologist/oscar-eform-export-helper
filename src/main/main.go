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
var sshHost, sshUser, sshPass, dbUser, dbPass, dbHost, dbName, dateFrom, dateTo, filePtr *string
var sshPort, fid *int

func main() {
	// Commandline flags
	sshHost = flag.String("host", "example.com", "The SSH host")
	sshPort = flag.Int("port", 22, "The port number")
	sshUser = flag.String("sshuser", "ssh-user", "ssh user")
	sshPass = flag.String("sshpass", "ssh-pass", "SSH Password")
	dbUser = flag.String("dbuser", "dbuser", "The db user")
	dbPass = flag.String("dbpass", "dbpass", "The db password")
	dbHost = flag.String("dbhost", "localhost:3306", "The db host")
	dbName = flag.String("dbname", "oscar", "The database name")
	dateFrom = flag.String("datefrom", "oscar", "The start date")
	dateTo = flag.String("dateto", "oscar", "The end date")
	fid = flag.Int("fid", 1, "The eform ID")
	filePtr = flag.String("file", "test.csv", "The csv file to process")
	flag.Parse()

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true

	g.SetManagerFunc(layout)

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
