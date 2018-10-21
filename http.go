package calendar

import (
	"fmt"
	"github.com/mvl-at/model"
	"net/http"
	"time"
)

const httpFormat = "20060102"

//Runs the http Server.
func run() {
	host := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	logger.Println("Listen on " + host)
	err := http.ListenAndServe(host, nil)

	if err != nil {
		errLogger.Fatalln(err.Error())
	}
}

//Registers all http routes.
func routes() {
	http.HandleFunc("/ical", events)
	http.HandleFunc("/pdf", pdf)
}

func events(rw http.ResponseWriter, r *http.Request) {
	events := eventsFromRange(r)
	convert := externalConvert
	if r.URL.Query().Get("int") == "true" {
		convert = musicianConvert
	}
	writeEvents(&events, rw, convert, threadType)
}

func pdf(rw http.ResponseWriter, r *http.Request) {
	events := eventsFromRange(r)
	note := r.URL.Query().Get("note")
	kpm, obm := fetchObmAndKpm()
	writeEventsTo(events, note, fmt.Sprintf("%s %s", obm.FirstName, obm.LastName), fmt.Sprintf("%s %s", kpm.FirstName, kpm.LastName), rw)
}

func eventsFromRange(r *http.Request) []*model.Event {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	last := r.URL.Query().Get("last")
	if last == "year" {
		from = time.Now().AddDate(-1, 0, 0).Format(httpFormat)
	}
	if last == "month" {
		from = time.Now().AddDate(0, -1, 0).Format(httpFormat)
	}
	if last == "week" {
		from = time.Now().AddDate(0, 0, -7).Format(httpFormat)
	}
	if last == "day" {
		from = time.Now().AddDate(0, 0, -1).Format(httpFormat)
	}
	events := make([]*model.Event, 0)
	fetchEvents(&events, from, to)
	return events
}
