package main

import (
	"fmt"
	"io/ioutil"
	// "log"
	"net/http"
	"regexp"
	"strings"

	// "github.com/gin-gonic/gin"
	"golang.org/x/net/html"
)

func main() {
	htmlCode, _ := readHtmlFromFile("./raumplan_e008.txt")

	line := extractTableLine(htmlCode)
	// log.Println(line)

	tableCells := parse(line)
	fmt.Println(len(tableCells))

	makePrettyOutput(tableCells)

	// r := gin.Default()
	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "pee pee pong",
	// 	})
	// })
	// r.Run("localhost:8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func readHtmlFromFile(fileName string) (string, error) {

	bs, err := ioutil.ReadFile(fileName)

	if err != nil {
		return "", err
	}

	return string(bs), nil
}

func extractTableLine(text string) (line string) {
	var woche1Table = regexp.MustCompile(`<table.*id='woche1'.*</table>`)
	line = woche1Table.FindString(text)
	return strings.Replace(line, "812px;", "812px;'", 1)
}

func makePrettyOutput(vals []string) (htmlOutput string) {

	tableHeaders := vals[0:7]

	for i, sadf := range tableHeaders {
		fmt.Printf("%v: %s\n", i, sadf)
	}

	// log.Println(tableHeaders)

	return "hello"

}

func parse(text string) (data []string) {

	tkn := html.NewTokenizer(strings.NewReader(text))
	var vals []string
	var insideTable bool
	var tableCellBuffer []string

	for {
		tt := tkn.Next()

		// log.Printf("currently on start tag token %s\n", tkn.Token().Data)
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
					// <table style='width:100%; max-width:812px; id='woche1'><tr><th style='width:50px;'>Uhrzeit</th><th>Montag</th><th>Dienstag</th><th>Mittwoch</th><th>Donnerstag</th><th>Freitag</th></tr><tr><td class='cent'>7:30 - 9:00</td><td></td><td></td><td></td><td><div>U Rechnerarchitektur I</div><span class='sml'><b>N.N.40 RA</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Rechnerarchitektur I</span></td><td><div>U SoI / HSC</div><span class='sml'><b>Wollschlaeger</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Systemorientierte Informatik und Hardware Software Codesign</span></td></tr><tr><td class='cent'>9:20 - 10:50</td><td><div>U Complexity Theory</div><span class='sml'><b>Krötzsch</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Complexity Theory</span></td><td><div>U HPC</div><span class='sml'><b>Nagel</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>High Performance Computing</span></td><td><div>V Principl. Dep. Syst.</div><span class='sml'><b>Fetzer</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Principles of Dependable Systems</span></td><td><div>U Informat.I/ ET</div><span class='sml'><b>N.N.10 ADS</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Informatik I für ET/MT/RES</span></td><td></td></tr><tr><td class='cent'>11:10 - 12:40</td><td><div>V Inf-Anw. Automation</div><span class='sml'><b>Wollschlaeger</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Informatik-Anwendungen in der Automation</span></td><td><div>U Betriebssysteme u. Sich.</div><span class='sml'><b>N.N.14 BSS</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Betriebssysteme und Sicherheit</span></td><td></td><td><div>U DB-Eng. Ü</div><span class='sml'><b>Lehner</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Datenbank-Engineering Übung </span></td><td><div>V Eng. Adapt. Mobile Apps</div><span class='sml'><b>Springer</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Engineering Adaptive Mobile Applications</span></td></tr><tr><td class='cent'>13:00 - 14:30</td><td><div>V Netzw. ind. Anw.</div><span class='sml'><b>Wollschlaeger</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Netzwerkmanagement in industriellen Anwendungen</span></td><td><div>U HS Techn. Datensch.</div><span class='sml'><b>Köpsell</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Hauptseminar Technischer Datenschutz</span></td><td><div>V PAofCS</div><span class='sml'><b>Nagel</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Performance Analysis of Computing Systems</span></td><td><div>U Principl. Dep. Syst.</div><span class='sml'><b>Fetzer</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Principles of Dependable Systems</span></td><td><div>U Eng. Adapt. Mobile Apps</div><span class='sml'><b>Springer</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Engineering Adaptive Mobile Applications</span></td></tr><tr><td class='cent'>14:50 - 16:20</td><td><div>V Information Retrieval</div><span class='sml'><b>Lehner</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Information Retrieval</span></td><td><div>U Microk.bas.</div><span class='sml'><b>Roitzsch</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Mikrokernbasierte Betriebssysteme</span></td><td><div>U Computergestützte Chirurgie</div><span class='sml'><b>Speidel</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Computer- und robotergestütze Chirurgie</span></td><td><div>U Form.Syst.</div><span class='sml'><b>N.N.33 FS</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Formale Systeme</span></td><td><div>U Informat.I/ ET</div><span class='sml'><b>N.N.13 ADS</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Informatik I für ET/MT/RES</span></td></tr><tr><td class='cent'>16:40 - 18:10</td><td><div>U Algorit.u.Daten</div><span class='sml'><b>N.N.03 AuD</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Algorithmen und Datenstrukturen</span></td><td><div>U Rechnerarchitektur I</div><span class='sml'><b>N.N.39 RA</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Rechnerarchitektur I</span></td><td><div>U SWT2</div><span class='sml'><b>N.N.17</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Softwaretechnologie II</span></td><td><div>U Info5 AuD</div><span class='sml'><b>N.N.INF</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Info5 AuD</span></td><td></td></tr><tr><td class='cent'>18:30 - 20:00</td><td></td><td></td><td></td><td><div>U Info5 AuD</div><span class='sml'><b>N.N.INF</b><span style='float:right'>INF</span></span><span class='more' style='display:none'>Info5 AuD</span></td><td></td></tr><tr><td class='cent'>20:20 - 21:50</td><td></td><td></td><td></td><td></td><td></td></tr><tr><td class='cent'>22:10 - 23:40</td><td></td><td></td><td></td><td></td><td></td></tr></table><p>Quelle: Datenbank Wintersemester</p>
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

func fetchTable() {

	url := "https://navigator.tu-dresden.de/raum/542100.2130"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("sec-ch-ua", "\"Google Chrome\";v=\"107\", \"Chromium\";v=\"107\", \"Not=A?Brand\";v=\"24\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-User", "?1")
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Cookie", "JSESSIONID=DDA7C652CE8C1185719A9C43CBA78397; loginToken=JyL8vUGd0x_HgDZAzyXN6")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
