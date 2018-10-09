package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"time"
)

func nextView(g *gocui.Gui, v *gocui.View) error {
	if v == nil || v.Name() == "side" {
		_, err := g.SetCurrentView("main")
		return err
	}
	_, err := g.SetCurrentView("side")
	return err
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	getLine(g, v)
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	getLine(g, v)
	return nil
}

func getLine(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	mainOutput(g, &l)
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("title", -1, -1, maxX, 3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorRed
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, "OSCAR eForm Export Tool Helper by Bell Eapen")
		fmt.Fprintln(v, "Valid Records: ", recordCount)
	}
	if _, err := g.SetView("main", 30, 4, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		message := "OSCAR Helper v 1.0.0"
		mainOutput(g, &message)
		if _, err := g.SetCurrentView("main"); err != nil {
			return err
		}
	}
	if v, err := g.SetView("side", -1, 4, 30, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		sideOutput(g)
	}
	return nil
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("side", gocui.KeyCtrlSpace, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyCtrlSpace, gocui.ModNone, nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("side", gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
		return err
	}
	return nil
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
