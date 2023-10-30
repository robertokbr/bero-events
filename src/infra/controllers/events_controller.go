package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/robertokbr/bero-events/src/domain/dtos"
)

type EventsController struct {
	jobs chan *dtos.EventDTO
}

func NewEventsController(jobs chan *dtos.EventDTO) *EventsController {
	return &EventsController{
		jobs: jobs,
	}
}

func (self *EventsController) Add(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data := dtos.EventDTO{}
	json.NewDecoder(r.Body).Decode(&data)

	self.jobs <- &data

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok": true}`))
}
