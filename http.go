package calendar

import (
	"fmt"
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

}
