package calendar

import (
	"github.com/mvl-at/model"
	"net/http"
)

const icalFormat = "20060102T150405"

func writeEvents(events *[]*model.Event, rw http.ResponseWriter) {
	eventChannel := make(chan *strBuffer, 1)
	for _, event := range *events {
		go eventToByte(event, eventChannel)
		(<-eventChannel).WriteTo(rw)
	}
}

func eventToByte(event *model.Event, eventChannel chan *strBuffer) {
	buf := &strBuffer{}
	buf.WriteFmt("BEGIN:EVENT")
	buf.WriteFmt("DTSTART:%s", event.Date.Format(icalFormat))
	buf.WriteFmt("SUMMARY:%s", event.Name)
	buf.WriteFmt("DESCRIPTION:%s", event.Note)
	buf.WriteFmt("LOCATION:%s", event.MusicianPlace)
	buf.WriteFmt("END:VEVENT")
	eventChannel <- buf
}
