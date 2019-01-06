package calendar

import (
	"github.com/mvl-at/model"
	"os"
	"testing"
	"time"
)

var testEvents = []*model.Event{
	{
		Name:          "Frühschoppen",
		Date:          time.Date(2018, 6, 17, 0, 0, 0, 0, time.Local),
		Uniform:       "MVL Polo und Lederhose",
		MusicianPlace: "Musikheim",
		Place:         "Musikheim",
		Time:          time.Date(1, 1, 1, 10, 0, 0, 0, time.Local),
		MusicianTime:  time.Date(0, 0, 0, 8, 45, 0, 0, time.Local),
		Note:          "Anschl. Wegräumen",
		Important:     true,
		Internal:      false},
	{
		Name:          "Marschmusikwertung",
		Date:          time.Date(2018, 9, 17, 0, 0, 0, 0, time.Local),
		Uniform:       "Uniform mit Hut",
		MusicianPlace: "Musikheim",
		Place:         "Spannberg",
		Time:          time.Date(1, 1, 1, 13, 0, 0, 0, time.Local),
		MusicianTime:  time.Date(0, 0, 0, 11, 0, 0, 0, time.Local),
		Important:     false,
		Internal:      false},

	{
		Name:          "Generalversammlung",
		Date:          time.Date(2018, 3, 17, 0, 0, 0, 0, time.Local),
		MusicianPlace: "Musikheim",
		MusicianTime:  time.Date(0, 0, 0, 20, 0, 0, 0, time.Local),
		Important:     true,
		Internal:      true}}

func TestOutput(t *testing.T) {
	out, err := os.OpenFile("test-plan.pdf", os.O_CREATE|os.O_RDWR, os.ModePerm)
	defer out.Close()
	if err != nil {
		t.Fatal(err.Error())
	}
	fpdf(testEvents, "Des is a notiz", "Willi Herok", "Markus Nentwich", out)
}

//func ExampleFpdf_MultiCell() {
//	pdf := gofpdf.New("P", "mm", "A4", "")
//	titleStr := "20000 Leagues Under the Seas"
//	pdf.SetTitle(titleStr, false)
//	pdf.SetAuthor("Jules Verne", false)
//	pdf.SetHeaderFunc(func() {
//		// Arial bold 15
//		pdf.SetFont("Arial", "B", 15)
//		// Calculate width of title and position
//		wd := pdf.GetStringWidth(titleStr) + 6
//		pdf.SetX((210 - wd) / 2)
//		// Colors of frame, background and text
//		pdf.SetDrawColor(0, 80, 180)
//		pdf.SetFillColor(230, 230, 0)
//		pdf.SetTextColor(220, 50, 50)
//		// Thickness of frame (1 mm)
//		pdf.SetLineWidth(1)
//		// Title
//		pdf.CellFormat(wd, 9, titleStr, "1", 1, "C", true, 0, "")
//		// Line break
//		pdf.Ln(10)
//	})
//	pdf.SetFooterFunc(func() {
//		// Position at 1.5 cm from bottom
//		pdf.SetY(-15)
//		// Arial italic 8
//		pdf.SetFont("Arial", "I", 8)
//		// Text color in gray
//		pdf.SetTextColor(128, 128, 128)
//		// Page number
//		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()),
//			"", 0, "C", false, 0, "")
//	})
//	chapterTitle := func(chapNum int, titleStr string) {
//		// 	// Arial 12
//		pdf.SetFont("Arial", "", 12)
//		// Background color
//		pdf.SetFillColor(200, 220, 255)
//		// Title
//		pdf.CellFormat(0, 6, fmt.Sprintf("Chapter %d : %s", chapNum, titleStr),
//			"", 1, "L", true, 0, "")
//		// Line break
//		pdf.Ln(4)
//	}
//	chapterBody := func(fileStr string) {
//		// Read text file
//		txtStr, err := ioutil.ReadFile(fileStr)
//		if err != nil {
//			pdf.SetError(err)
//		}
//		// Times 12
//		pdf.SetFont("Times", "", 12)
//		// Output justified text
//		pdf.MultiCell(0, 5, string(txtStr), "", "", false)
//		// Line break
//		pdf.Ln(-1)
//		// Mention in italics
//		pdf.SetFont("", "I", 0)
//		pdf.Cell(0, 5, "(end of excerpt)")
//	}
//printChapter := func(chapNum int, titleStr, fileStr string) {
//	pdf.AddPage()
//	chapterTitle(chapNum, titleStr)
//	chapterBody(fileStr)
//}
//printChapter(1, "A RUNAWAY REEF", example.TextFile("20k_c1.txt"))
//printChapter(2, "THE PROS AND CONS", example.TextFile("20k_c2.txt"))
//fileStr := example.Filename("Fpdf_MultiCell")
//err := pdf.OutputFileAndClose(fileStr)
//example.Summary(err, fileStr)
// Output:
// Successfully generated pdf/Fpdf_MultiCell.pdf
//}
