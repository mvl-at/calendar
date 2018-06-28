package calendar

import (
	"encoding/json"
	"fmt"
	"github.com/mvl-at/model"
	"io/ioutil"
	"net/http"
)

func fetchEvents(events *[]*model.Event, from string, to string) {
	var jsonData []byte
	resp, err := http.Get(fmt.Sprintf("http://%s/eventsrange?from=%s&to=%s", conf.RestHost, from, to))

	if err != nil {
		errLogger.Println(err.Error())
		return
	}
	if resp.Body != nil {
		jsonData, _ = ioutil.ReadAll(resp.Body)
	} else {
		return
	}
	json.Unmarshal(jsonData, events)
}
