package calendar

import (
	"fmt"
	"github.com/mvl-at/model"
	"time"
)

var weekDays = map[time.Weekday]string{
	time.Sunday:    "Sonntag",
	time.Monday:    "Montag",
	time.Tuesday:   "Dienstag",
	time.Wednesday: "Mittwoch",
	time.Thursday:  "Donnerstag",
	time.Friday:    "Freitag",
	time.Saturday:  "Samstag"}

var months = map[time.Month]string{
	time.January:   "Jänner",
	time.February:  "Feber",
	time.March:     "März",
	time.April:     "April",
	time.May:       "Mai",
	time.June:      "Juni",
	time.July:      "Juli",
	time.August:    "August",
	time.September: "September",
	time.October:   "Oktober",
	time.November:  "November",
	time.December:  "Dezember"}

type RenderedEvent struct {
	Date    string
	Name    string
	Begin   string
	Venue   string
	Note    string
	Uniform string

	HasDate    bool
	HasName    bool
	HasBegin   bool
	HasVenue   bool
	HasNote    bool
	HasUniform bool
}

func renderEvent(event *model.Event) *RenderedEvent {
	re := &RenderedEvent{
		Name:    event.Name,
		Note:    event.Note,
		Uniform: event.Uniform,
		Date:    fmt.Sprintf(event.Date.Format("%s, 02.01.06"), weekDays[event.Date.Weekday()]),
		Begin:   event.Place + ", " + event.Time.Format("15:04 Uhr"),
		Venue:   event.MusicianPlace + ", " + event.MusicianTime.Format("15:04 Uhr")}
	re.HasBegin = !event.Internal
	re.HasDate = re.Date != ""
	re.HasName = re.Name != ""
	re.HasNote = re.Note != ""
	re.HasUniform = re.Uniform != ""
	re.HasVenue = re.Venue != ""
	return re
}

func renderAllEvents(events []*model.Event) []*RenderedEvent {
	renderedEvents := make([]*RenderedEvent, len(events))
	for i := range events {
		renderedEvents[i] = renderEvent(events[i])
	}
	return renderedEvents
}
