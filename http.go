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
	http.HandleFunc("/events", events)
}

func events(rw http.ResponseWriter, r *http.Request) {
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
	convert := externalConvert
	if r.URL.Query().Get("int") == "true" {
		convert = musicianConvert
	}
	writeEvents(&events, rw, convert, threadType)
}
