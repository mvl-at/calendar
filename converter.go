package calendar

import (
	"github.com/mvl-at/model"
	"net/http"
	"runtime"
)

const icalFormat = "20060102T150405"

func writeEvents(events *[]*model.Event, rw http.ResponseWriter) {
	threads := conf.Threads
	if threads < 1 {
		threads = runtime.NumCPU()
	}
	if len(*events) < threads {
		threads = len(*events)
	}
	eventChannel := make(chan *strBuffer, threads)
	for i := 0; i < threads; i++ {
		go eventToByte((*events)[i], eventChannel)
	}

	for i := 0; i < len(*events); i++ {
		(<-eventChannel).WriteTo(rw)
		if i < len(*events)-threads {
			go eventToByte((*events)[i+threads], eventChannel)
		}
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
