package usecases

import (
	"fmt"

	"github.com/robertokbr/events-worker/src/domain/dtos"
)

func EventsWorker(jobs <-chan *dtos.EventDTO, jobsReturn chan JobsReturn, claimBadgesUsecase *ClaimBadgesUsecase, claimRewardsUsecase *ClaimRewardsUsecase, id int) {
	for job := range jobs {
		switch job.Event {
		case dtos.LOOT_COLLECTED:
			go claimBadgesUsecase.Execute(job.UserID)
			break
		case dtos.REWARD_COLLECTED:
			go claimRewardsUsecase.Execute(job.UserID)
			break
		default:
			fmt.Printf("User %d did something", *&job.UserID)
		}
	}
}
