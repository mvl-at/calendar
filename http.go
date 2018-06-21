package calendar

import (
	"fmt"
	"github.com/mvl-at/model"
	"net/http"
)

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
	events := make([]*model.Event, 0)
	fetchEvents(&events)
	writeEvents(&events, rw)
}
