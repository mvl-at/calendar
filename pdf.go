package calendar

import (
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/mvl-at/model"
	"html/template"
	"io"
	"strconv"
	"strings"
)

const (
	eventTemplate = "events.html"
)

type EventListInfo struct {
	Events    []*RenderedEvent
	Obm       string
	Kpm       string
	Note      string
	HasNote   bool
	Range     string
	MarchMode bool
	Marches   []string
}

func writeEventsTo(events []*model.Event, note string, obm string, kpm string, writer io.Writer) {
	t, _ := template.ParseFiles(eventTemplate)
	rangeString := "für gute Märsche"
	marchMode := false

	if len(events) > 0 {
		firstYear := events[0].Date.Year()
		lastYear := events[len(events)-1].Date.Year()
		firstMonth := events[0].Date.Month()
		lastMonth := events[len(events)-1].Date.Month()

		if firstYear == lastYear {
			if firstMonth == lastMonth {
				rangeString = months[firstMonth] + " " + strconv.Itoa(firstYear)
			} else {
				rangeString = months[firstMonth] + " bis " + months[lastMonth] + " " + strconv.Itoa(firstYear)
			}
		} else {
			rangeString = months[firstMonth] + " " + strconv.Itoa(firstYear) + " bis " + months[lastMonth] + " " + strconv.Itoa(lastYear)
		}
	} else {
		if strings.ToLower(note) == strings.ToLower(conf.King) {
			rangeString = "gute Märsche"
			marchMode = true
		} else {
			note = conf.King
		}
	}

	data := EventListInfo{Events: renderAllEvents(events), Obm: obm, Kpm: kpm, Note: note, HasNote: note != "", Range: rangeString, MarchMode: marchMode, Marches: conf.Marches}

	pdfg, _ := wkhtmltopdf.NewPDFGenerator()
	pdfg.Title.Set("Termine")
	pdfg.Dpi.Set(600)
	pr, pw := io.Pipe()
	defer pr.Close()
	go func() {
		defer pw.Close()
		t.Execute(pw, data)
	}()
	page := wkhtmltopdf.NewPageReader(pr)
	pdfg.AddPage(page)
	err := pdfg.Create()
	if err != nil {
		errLogger.Println(err.Error())
		return
	}
	if err != nil {
		errLogger.Println(err.Error())
		return
	}
	io.Copy(writer, pdfg.Buffer())
}
