package calendar

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"github.com/mvl-at/model"
	"io"
	"os"
	"sort"
)

const (
	headerSize  = 18
	stdSize     = 12
	smallSize   = 9
	eventMargin = 4
	infoMargin  = 4
	infoXMargin = 32
	boxMarginY  = 22
)

type RenderedEventField struct {
	FieldDescription func(event *RenderedEvent) (str string, bold bool)
	FieldValue       func(event *RenderedEvent) (str string, bold bool)
	Render           func(event *RenderedEvent) bool
}

func fpdf(events []*model.Event, note string, author string, writer io.Writer) {
	sort.Slice(events, func(i, j int) bool {
		return events[i].Date.Unix() < events[j].Date.Unix()
	})
	pdf := gofpdf.New("P", "mm", "A4", "")
	defer pdf.Close()
	pdf.SetAuthor(author, true)
	pdf.SetCreator("Musikverein Leopoldsdorf Webseite", false)
	pdf.SetAutoPageBreak(true, 4)
	pdf.SetFont("Arial", "", stdSize)
	pdf.AddPage()
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	center := func(text string, link string) float64 {
		pageWidth, _ := pdf.GetPageSize()
		textWidth := pdf.GetStringWidth(text) + 2
		pdf.SetX((pageWidth - textWidth) / 2)
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
		title := conf.Name
		beginY := pdf.GetY()
		widths := make([]float64, 4)
		widths[0] = center(title, "")
		pdf.SetFontSize(smallSize)
		widths[1] = center(conf.HomePage, "")
		widths[2] = center(conf.ZVR, "")
		widths[3] = center(conf.Address, "")
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
		mvlWidth := ldfWidth
		mvlHeight := mvlWidth / mvlLogo.Width() * mvlLogo.Height()
		pdf.Image("ldf", ((pageWidth-boxWidth)/2-ldfWidth)/2, beginY, ldfWidth, boxHeight, false, "", 0, "")
		pdf.Image("mvl", (pageWidth-(pageWidth-boxWidth)/4)-mvlWidth/2, beginY+(boxHeight-mvlHeight)/2, mvlWidth, mvlHeight, false, "", 0, "")
		pdf.Rect((pageWidth-boxWidth)/2, beginY, boxWidth, boxHeight, "D")
	}

	drawEventLine := func(text string, bold bool, line int) float64 {
		text = tr(text)
		pdf.SetFontSize(stdSize)
		oldX := pdf.GetX()
		pdf.SetX(oldX + eventMargin)
		textWidth := pdf.GetStringWidth(text)
		fontSize, _ := pdf.GetFontSize()
		if bold {
			pdf.SetFont("", "B", 12)
		}
		pdf.Write(fontSize/2, text)
		pdf.SetX(oldX)
		if bold {
			pdf.SetFont("", "", 12)
		}
		return textWidth
	}

	drawEvent := func(event *RenderedEvent, renderer *RenderedEventField, maxWidth float64) {
		if !renderer.Render(event) {
			return
		}
		oldX := pdf.GetX()
		str, bold := renderer.FieldDescription(event)
		drawEventLine(str, bold, 0)
		pdf.SetX(oldX + maxWidth + eventMargin)
		str, bold = renderer.FieldValue(event)
		drawEventLine(str, bold, 1)
		fontSize, _ := pdf.GetFontSize()
		pdf.Ln(fontSize / 2)
		pdf.SetX(oldX)
	}

	evs := renderAllEvents(events)

	drawEvents := func() {
		pdf.SetFont("", "B", headerSize)
		pdf.SetY(pdf.GetY() + eventMargin)
		title := rangeString(events, note)
		pdf.SetTitle(title, true)
		center(tr(title), "")
		fontSize, _ := pdf.GetFontSize()
		pdf.Ln(fontSize / 2)
		pdf.SetFont("", "", stdSize)
		widths := make([]float64, 0)
		oldY := pdf.GetY()
		fields := []*RenderedEventField{
			{
				FieldDescription: func(event *RenderedEvent) (str string, bold bool) {
					return event.Date, true
				},
				FieldValue: func(event *RenderedEvent) (str string, bold bool) {
					return event.Name, true
				},
				Render: func(event *RenderedEvent) bool {
					return true
				},
			},
			{
				FieldDescription: func(event *RenderedEvent) (str string, bold bool) {
					return "Treffpunkt:", false
				},
				FieldValue: func(event *RenderedEvent) (str string, bold bool) {
					return event.Venue, false
				},
				Render: func(event *RenderedEvent) bool {
					return event.HasVenue
				},
			},
			{
				FieldDescription: func(event *RenderedEvent) (str string, bold bool) {
					return "Beginn:", false
				},
				FieldValue: func(event *RenderedEvent) (str string, bold bool) {
					return event.Begin, false
				},
				Render: func(event *RenderedEvent) bool {
					return event.HasBegin
				},
			},
			{
				FieldDescription: func(event *RenderedEvent) (str string, bold bool) {
					return "Adjustierung:", false
				},
				FieldValue: func(event *RenderedEvent) (str string, bold bool) {
					return event.Uniform, false
				},
				Render: func(event *RenderedEvent) bool {
					return event.HasUniform
				},
			},
			{
				FieldDescription: func(event *RenderedEvent) (str string, bold bool) {
					return "Ende:", false
				},
				FieldValue: func(event *RenderedEvent) (str string, bold bool) {
					return event.End, false
				},
				Render: func(event *RenderedEvent) bool {
					return event.HasEnd
				},
			},
			{
				FieldDescription: func(event *RenderedEvent) (str string, bold bool) {
					return "Notiz:", false
				},
				FieldValue: func(event *RenderedEvent) (str string, bold bool) {
					return event.Note, false
				},
				Render: func(event *RenderedEvent) bool {
					return event.HasNote
				},
			},
		}
		for _, event := range evs {
			for _, field := range fields {
				if field.Render(event) {
					str, _ := field.FieldDescription(event)
					widths = append(widths, pdf.GetStringWidth(str))
				}
			}
		}
		pdf.SetY(oldY)
		gr := greatest(&widths)

		for i, v := range evs {
			if i != 0 {
				fontSize, _ := pdf.GetFontSize()
				pdf.Ln(fontSize / 2)
			}
			for _, field := range fields {
				drawEvent(v, field, gr)
			}

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
		pdf.SetX(infoXMargin)
		html.Write(smallSize, personString(conf.Obm, true, "Obmann"))
		pdf.Ln(smallSize / 2)
		pdf.SetX(infoXMargin)
		html.Write(smallSize, personString(conf.Kpm, true, "Kapellmeister"))
		pdf.Ln(smallSize / 2)
		pdf.SetY(pdf.GetY() + eventMargin)
		pdf.Line(pdf.GetX()+eventMargin, pdf.GetY(), pageWidth-pdf.GetX()-eventMargin, pdf.GetY())
	}

	noteBox := func() {
		if len(note) > 0 {
			width, height := pdf.GetPageSize()
			pdf.SetY(height - boxMarginY)
			pdf.SetFontSize(smallSize)
			pdf.SetX(infoXMargin)
			pdf.MultiCell(width-(infoXMargin*2), smallSize/1.5, tr(note), "1", "c", false)
		}
	}

	header()
	infoHeader()
	pdf.SetX(0)
	drawEvents()
	noteBox()

	if err := pdf.Output(writer); err != nil {
		errLogger.Println(err.Error())
	}
}
