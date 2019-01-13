package calendar

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"github.com/mvl-at/model"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

func rangeString(events []*model.Event, note string) string {
	rangeString := "Terminplan für gute Märsche"
	if strings.ToLower(note) == strings.ToLower(conf.King) {
		rangeString = "Gute Märsche"
	}
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
		rangeString = "Terminplan " + rangeString
	}
	return rangeString
}

const (
	headerSize            = 18
	stdSize               = 12
	smallSize             = 9
	eventMargin           = 4
	infoMargin            = 4
	infoXMargin           = 32
)

func fpdf(events []*model.Event, note string, writer io.Writer) {
	sort.Slice(events, func(i, j int) bool {
		return events[i].Date.Unix() < events[j].Date.Unix()
	})
	pdf := gofpdf.New("P", "mm", "A4", "")
	defer pdf.Close()
	pdf.SetAuthor("Karl Steinscheisser", true)
	pdf.SetAutoPageBreak(true, 10)
	pdf.SetFont("Arial", "", stdSize)
	pdf.AddPage()
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	center := func(text string, link string) float64 {
		pageWidth, _ := pdf.GetPageSize()
		textWidth := pdf.GetStringWidth(text) + 2
		pdf.SetX((pageWidth - textWidth) / 2)
		//pdf.Write(0, text)
		fontSize, _ := pdf.GetFontSize()
		isLink := 0
		if link != "" {
			isLink = 1
		}
		pdf.CellFormat(textWidth+4, fontSize/2, text, "0", 1, "C", false, isLink, link)
		return textWidth
	}

	greatest := func(sl *[]float64) float64 {
		gr := (*sl)[0]
		for _, v := range *sl {
			if v > gr {
				gr = v
			}
		}
		return gr
	}
	header := func() {
		pdf.SetFont("", "B", headerSize)
		title := "Musikverein Leopoldsdorf/M."
		beginY := pdf.GetY()
		widths := make([]float64, 4)
		widths[0] = center(title, "")
		pdf.SetFontSize(smallSize)
		widths[1] = center("www.mvl.at", "")
		widths[2] = center("ZVR - Zahl: 091786949", "")
		widths[3] = center("A-2285 Leopoldsdorf/M. Kempfendorf 2", "")
		pageWidth, _ := pdf.GetPageSize()
		boxWidth := greatest(&widths) + 10
		boxHeight := pdf.GetY() - beginY
		ldfF, err := os.Open("ldf.png")
		defer ldfF.Close()
		mvlF, err := os.Open("mvl.png")
		defer mvlF.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
		ldfLogo := pdf.RegisterImageReader("ldf", "png", ldfF)
		mvlLogo := pdf.RegisterImageReader("mvl", "png", mvlF)
		ldfWidth := (boxHeight / ldfLogo.Height()) * ldfLogo.Width()
		mvlWidth := ldfWidth * 1.5
		mvlHeight := mvlWidth / mvlLogo.Width() * mvlLogo.Height()
		pdf.Image("ldf", ((pageWidth-boxWidth)/2-ldfWidth)/2, beginY, ldfWidth, boxHeight, false, "", 0, "")
		pdf.Image("mvl", (pageWidth-(pageWidth-boxWidth)/4)-mvlWidth/2, beginY+(boxHeight-mvlHeight)/2, mvlWidth, mvlHeight, false, "", 0, "")
		pdf.Rect((pageWidth-boxWidth)/2, beginY, boxWidth, boxHeight, "D")
	}

	drawEventLine := func(text string, bold bool) float64 {
		text = tr(text)
		pdf.SetFontSize(stdSize)
		oldX := pdf.GetX()
		pdf.SetX(oldX + eventMargin)
		textWidth := pdf.GetStringWidth(text)
		fontSize, _ := pdf.GetFontSize()
		if bold {
			pdf.SetFont("", "B", 12)
		}
		pdf.CellFormat(textWidth+4, fontSize/2, text, "0", 1, "C", false, 0, "")
		pdf.SetX(oldX)
		if bold {
			pdf.SetFont("", "", 12)
		}
		return textWidth
	}

	drawEvent := func(event *RenderedEvent, widths *[]float64, descriptors bool) {
		if descriptors {
			*widths = append(*widths, drawEventLine(event.Date, true))
		} else {
			drawEventLine(event.Name, true)
		}
		if event.HasVenue {
			if descriptors {
				*widths = append(*widths, drawEventLine("Treffpunkt:", false))
			} else {
				drawEventLine(event.Venue, false)
			}
		}
		if event.HasBegin {
			if descriptors {
				*widths = append(*widths, drawEventLine("Beginn:", false))
			} else {
				drawEventLine(event.Begin, false)
			}
		}
		if event.HasUniform {
			if descriptors {
				*widths = append(*widths, drawEventLine("Adjustierung:", false))
			} else {
				drawEventLine(event.Uniform, false)
			}
		}
		if event.HasNote {
			if descriptors {
				*widths = append(*widths, drawEventLine("Notiz:", false))
			} else {
				drawEventLine(event.Note, false)
			}
		}
	}

	evs := renderAllEvents(events)

	drawEvents := func() {
		pdf.SetFont("", "B", headerSize)
		pdf.SetY(pdf.GetY() + eventMargin)
		center(tr(rangeString(events, note)), "")
		pdf.SetFont("", "", stdSize)
		widths := make([]float64, 0)
		oldY := pdf.GetY()
		for _, v := range evs {
			pdf.SetY(pdf.GetY() + eventMargin)
			drawEvent(v, &widths, true)
		}
		pdf.SetY(oldY)
		gr := greatest(&widths)
		oldX := pdf.GetX()

		for _, v := range evs {
			pdf.SetY(pdf.GetY() + eventMargin)
			pdf.SetX(oldX + gr + eventMargin)
			drawEvent(v, &widths, false)
		}
	}

	personString := func(person Person, bold bool, role string) string {
		boldStart := "<b>"
		boldEnd := "</b>"
		mailStart := "<a href=\"mailto:" + person.Email + "\">"
		mailEnd := "</a>"
		if !bold {
			boldEnd = ""
			boldStart = ""
			mailStart = ""
			mailEnd = ""
		}
		return fmt.Sprintf("%s%s:%s %s, %s, %s, %s%s%s", boldStart, role, boldEnd, person.Name, person.Address, person.Telephone, mailStart, person.Email, mailEnd)
	}

	infoHeader := func() {
		pdf.SetY(pdf.GetY() + infoMargin)
		pdf.SetFontSize(smallSize)
		html := pdf.HTMLBasicNew()
		pageWidth, _ := pdf.GetPageSize()
		//pdf.SetX((pageWidth - pdf.GetStringWidth(personString(conf.Obm, false, "Obmann"))) / 2)
		pdf.SetX(infoXMargin)
		html.Write(smallSize, personString(conf.Obm, true, "Obmann"))
		pdf.Ln(smallSize / 2)
		//pdf.SetX((pageWidth - pdf.GetStringWidth(personString(conf.Kpm, false, "Kapellmeister"))) / 2)
		pdf.SetX(infoXMargin)
		html.Write(smallSize, personString(conf.Kpm, true, "Kapellmeister"))
		pdf.Ln(smallSize / 2)
		pdf.SetY(pdf.GetY() + eventMargin)
		pdf.Line(pdf.GetX()+eventMargin, pdf.GetY(), pageWidth-pdf.GetX()-eventMargin, pdf.GetY())

	}

	header()
	infoHeader()
	pdf.SetX(0)
	drawEvents()

	fmt.Println(pdf.Output(writer))
}
