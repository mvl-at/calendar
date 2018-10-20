package calendar

import (
	"encoding/json"
	"fmt"
	"github.com/mvl-at/model"
	"io/ioutil"
	"net/http"
	"sort"
	"time"
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
	eventsVal := *events
	sort.Slice(eventsVal, func(i, j int) bool {
		iD := eventsVal[i].Date
		iT := eventsVal[i].Time
		iTime := time.Date(iD.Year(), iD.Month(), iD.Day(), iT.Hour(), iT.Minute(), 0, 0, time.Local)
		jD := eventsVal[j].Date
		jT := eventsVal[j].Time
		jTime := time.Date(jD.Year(), jD.Month(), jD.Day(), jT.Hour(), jT.Minute(), 0, 0, time.Local)
		return iTime.Unix() < jTime.Unix()
	})
}
