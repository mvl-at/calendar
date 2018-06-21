package calendar

import (
	"github.com/mvl-at/model"
	"net/http"
	"testing"
	"time"
)

func BenchmarkThreading(b *testing.B) {
	e := &model.Event{Name: "Frühschoppen", Date: time.Date(2018, 6, 17, 9, 0, 0, 0, time.Local)}
	events := make([]*model.Event, 1000000)
	for i := range events {
		events[i] = e
	}
	conf = &Configuration{Threads: 4, RestHost: "127.0.0.1:8071", Host: "127.0.0.1", Port: 8071}
	b.ResetTimer()
	writeEvents(&events, dummyWriter{})
}

type dummyWriter struct {
}

func (d dummyWriter) Header() http.Header        { return nil }
func (d dummyWriter) Write([]byte) (int, error)  { return 0, nil }
func (d dummyWriter) WriteHeader(statusCode int) {}
