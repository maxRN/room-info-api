package main

import (
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// fixMalformedHTML fixes the broken HTML returned by the TU servers
func fixMalformedHTML(text string) (line string) {
	var woche1Table = regexp.MustCompile(`<table.*id='woche1'.*</table>`)
	line = woche1Table.FindString(text)
	return strings.Replace(line, "812px;", "812px;'", 1)
}

func parseIntoRoomInfo(vals lecturePlan) (info RoomInfo) {

	ds1 := removeNewLines(vals[7:12])
	ds2 := removeNewLines(vals[13:18])
	ds3 := removeNewLines(vals[19:24])
	ds4 := removeNewLines(vals[25:30])
	ds5 := removeNewLines(vals[31:36])
	ds6 := removeNewLines(vals[37:42])
	ds7 := removeNewLines(vals[43:48])
	ds8 := removeNewLines(vals[49:54])
	ds9 := removeNewLines(vals[55:60])
	ds10 := []string{"", "", "", "", ""}

	roomInfo := RoomInfo{
		Ds1:  ds1,
		Ds2:  ds2,
		Ds3:  ds3,
		Ds4:  ds4,
		Ds5:  ds5,
		Ds6:  ds6,
		Ds7:  ds7,
		Ds8:  ds8,
		Ds9:  ds9,
		Ds10: ds10,
	}

	return roomInfo

}

type lecturePlan = []string

func parseIntoRoom(rr rawRoom) (r Room) {
	planInfo := extractLecturePlanInfo(rr.WebPage)
	ri := parseIntoRoomInfo(planInfo)
	r.Name = rr.Name
	r.Building = rr.Building
	r.Ds1 = ri.Ds1
	r.Ds2 = ri.Ds2
	r.Ds3 = ri.Ds3
	r.Ds4 = ri.Ds4
	r.Ds5 = ri.Ds5
	r.Ds6 = ri.Ds6
	r.Ds7 = ri.Ds7
	r.Ds8 = ri.Ds8
	r.Ds9 = ri.Ds9
	r.Ds10 = ri.Ds10

	return r
}

// parses the raw HTML scraped from the TU website into the lecture plan
func extractLecturePlanInfo(htmlCode string) (plan lecturePlan) {
	htmlCode = fixMalformedHTML(htmlCode)

	tkn := html.NewTokenizer(strings.NewReader(htmlCode))
	var vals []string
	var insideTable bool
	var tableCellBuffer []string

	for {
		tt := tkn.Next()

		switch {
		case tt == html.ErrorToken:
			return vals
		case tt == html.StartTagToken:
			t := tkn.Token()

			attributes := t.Attr
			for _, attribute := range attributes {
				if attribute.Key == "id" && attribute.Val == "woche1" {
					// this is the table we are looking for
					// then return everything inside
					insideTable = true
				}
			}
		case tt == html.TextToken:
			t := tkn.Token()
			if insideTable {
				tableCellBuffer = append(tableCellBuffer, t.Data)

			}
		case tt == html.EndTagToken:
			t := tkn.Token()
			if t.Data == "td" || t.Data == "th" {
				vals = append(vals, strings.Join(tableCellBuffer, "\n"))
				tableCellBuffer = nil
			}
			if insideTable && t.Data == "table" {
				insideTable = false
			}
		}
	}
}

func removeNewLines(vals []string) (noLines []string) {
	longString := strings.Join(vals, ";")
	noNewLines := strings.ReplaceAll(longString, "\n", " ")
	return strings.Split(noNewLines, ";")
}
