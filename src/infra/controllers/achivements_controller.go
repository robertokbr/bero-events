package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/robertokbr/bero-events/src/domain/dtos"
)

type AchievementsController struct {
	jobs chan *dtos.CheckAchievementsDTO
}

func NewAchievementsController(jobs chan *dtos.CheckAchievementsDTO) *AchievementsController {
	return &AchievementsController{
		jobs: jobs,
	}
}

func (self *AchievementsController) CheckAchievements(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data := dtos.CheckAchievementsDTO{}
	json.NewDecoder(r.Body).Decode(&data)

	self.jobs <- &data

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok": true}`))
}
