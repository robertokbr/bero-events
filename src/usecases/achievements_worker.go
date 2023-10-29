package usecases

import (
	"fmt"

	"github.com/robertokbr/bero-events/src/domain/dtos"
)

func AchievementsWorker(jobs <-chan *dtos.CheckAchievementsDTO, jobsReturn chan JobsReturn, badgesUsecase *BadgesUsecase, id int) {
	for job := range jobs {
		switch job.Event {
		case dtos.LOOT_COLLECTED:
			badgesUsecase.Execute(job.UserID)
			break
		case dtos.DAILY_INTERACTION_MADE:
			fmt.Printf("User %d made a daily interaction", *&job.UserID)
			break
		default:
			fmt.Printf("User %d did something", *&job.UserID)
		}
	}
}
