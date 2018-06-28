package calendar

import (
	"github.com/mvl-at/model"
	"net/http"
	"runtime"
	"time"
)

//date format string which is used in the ical format
const icalFormat = "20060102T150405"

//main wrapper for event converter
func writeEvents(events *[]*model.Event, rw http.ResponseWriter, convertEvent convertEvent, control threadControl) {
	threads := conf.Threads
	if threads < 1 {
		threads = runtime.NumCPU()
	}
	if len(*events) < threads {
		threads = len(*events)
	}
	writeHeader(rw)
	control(events, rw, threads, convertEvent)
	writeFooter(rw)
}

func writeHeader(rw http.ResponseWriter) {
	buf := &strBuffer{}
	buf.WriteFmt("BEGIN:VCALENDAR")
	buf.WriteFmt("VERSION:2.0")
	buf.WriteTo(rw)
}

func writeFooter(rw http.ResponseWriter) {
	buf := &strBuffer{}
	buf.WriteFmt("END:VCALENDAR")
	buf.WriteTo(rw)
}

//function to write a event into a string buffer
type convertEvent func(event *model.Event, buffer *strBuffer)

//function which defines how threads are gonna created and destroyed while converting events
type threadControl func(events *[]*model.Event, rw http.ResponseWriter, threads int, convertEvent convertEvent)

//converts events of the view of a musician
func musicianConvert(event *model.Event, buffer *strBuffer) {
	buffer.WriteFmt("BEGIN:VEVENT")
	t := event.MusicianTime
	buffer.WriteFmt("DTSTART:%s", event.Date.Add(time.Duration(t.Hour())*time.Hour+time.Duration(t.Minute())*time.Minute+time.Duration(t.Second())*time.Second).Format(icalFormat))
	buffer.WriteFmt("SUMMARY:%s", event.Name)
	buffer.WriteFmt("DESCRIPTION:%s (%s)", event.Note, event.Uniform)
	buffer.WriteFmt("LOCATION:%s", event.MusicianPlace)
	buffer.WriteFmt("END:VEVENT")
}

//converts events of the view of a non-musician
func externalConvert(event *model.Event, buffer *strBuffer) {
	if !event.Internal {
		buffer.WriteFmt("BEGIN:VEVENT")
		t := event.Time
		buffer.WriteFmt("DTSTART:%s", event.Date.Add(time.Duration(t.Hour())*time.Hour+time.Duration(t.Minute())*time.Minute+time.Duration(t.Second())*time.Second).Format(icalFormat))
		buffer.WriteFmt("SUMMARY:%s", event.Name)
		buffer.WriteFmt("LOCATION:%s", event.Place)
		buffer.WriteFmt("END:VEVENT")
	}
}

//destroys thread if done and creates a new one
func dynamicThreads(events *[]*model.Event, rw http.ResponseWriter, threads int, convertEvent convertEvent) {
	eventChannel := make(chan *strBuffer, threads)
	for i := 0; i < threads; i++ {
		go singleIterate((*events)[i], convertEvent, eventChannel)
	}
	for i := 0; i < len(*events); i++ {
		(<-eventChannel).WriteTo(rw)
		if i < len(*events)-threads {
			go singleIterate((*events)[i], convertEvent, eventChannel)
		}
	}
}

//uses multiple buffers to convert events
func singleIterate(event *model.Event, convertEvent convertEvent, buffer chan *strBuffer) {
	buf := &strBuffer{}
	convertEvent(event, buf)
	buffer <- buf
}

//static thread creation
func staticThreads(events *[]*model.Event, rw http.ResponseWriter, threads int, convertEvent convertEvent) {
	startIndex := 0
	itemsPerThread := len(*events) / threads
	endIndex := itemsPerThread
	bufferChannel := make(chan *strBuffer)
	for i := 0; i < threads; i++ {
		eventRange := (*events)[startIndex:endIndex]
		if i == threads-1 {
			eventRange = (*events)[startIndex:]
		}
		go multiIterate(&eventRange, convertEvent, bufferChannel)
		//oldStart := startIndex
		startIndex = endIndex
		endIndex = startIndex + itemsPerThread
	}
	for i := 0; i < threads; i++ {
		(<-bufferChannel).WriteTo(rw)
	}
}

//uses one buffer to convert multiple events
func multiIterate(events *[]*model.Event, convertEvent convertEvent, buffer chan *strBuffer) {
	buf := &strBuffer{}
	for _, event := range *events {
		convertEvent(event, buf)
	}
	buffer <- buf
}

//converts all events in the main thread
func mainThread(events *[]*model.Event, rw http.ResponseWriter, _ int, convertEvent convertEvent) {
	buf := &strBuffer{}
	for _, event := range *events {
		convertEvent(event, buf)
	}
	buf.WriteTo(rw)
}
