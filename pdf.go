package calendar

import (
	"fmt"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/mvl-at/model"
	"html/template"
	"io"
	"os"
	"strconv"
)

const (
	eventTemplate = "events.html"
)

type EventListInfo struct {
	Events  []*RenderedEvent
	Obm     string
	Kpm     string
	Note    string
	HasNote bool
	Range   string
}

func writeEventsTo(events []*model.Event, note string, writer io.Writer) {
	t, _ := template.ParseFiles(eventTemplate)
	rangeString := ""
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

	data := EventListInfo{Events: renderAllEvents(events), Obm: "Wilhelm Herok", Kpm: "Markus Nentwich", Note: note, HasNote: note != "", Range: rangeString}

	pdfg, _ := wkhtmltopdf.NewPDFGenerator()
	pdfg.Title.Set("Termine")
	pdfg.Dpi.Set(6000)
	wd, _ := os.Getwd()
	fmt.Println(wd)
	err := os.Symlink("res", os.TempDir()+"/res")
	if err != nil {
		fmt.Println(err.Error())
	}
	pr, pw := io.Pipe()
	defer pr.Close()
	go func() {
		defer pw.Close()
		t.Execute(pw, data)
	}()
	page := wkhtmltopdf.NewPageReader(pr)
	pdfg.AddPage(page)
	err = pdfg.Create()
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
