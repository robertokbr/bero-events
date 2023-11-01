package usecases

import (
	"github.com/robertokbr/events-worker/src/domain/dtos"
	"github.com/robertokbr/events-worker/src/domain/interfaces"
	"github.com/robertokbr/events-worker/src/infra/database/repositories"
)

type JobsReturn struct {
	Job *dtos.EventDTO
	Err error
}

type EventsWorkerManager struct {
	repository   *repositories.MySqlRepository
	mailProvider interfaces.MailProvider
	jobs         chan *dtos.EventDTO
}

func NewEventsWorkerManager(jobs chan *dtos.EventDTO, repository *repositories.MySqlRepository, mailProvider interfaces.MailProvider) *EventsWorkerManager {
	return &EventsWorkerManager{
		jobs:         jobs,
		repository:   repository,
		mailProvider: mailProvider,
	}
}

func (self *EventsWorkerManager) Start(numberOfThreads int) {
	jobsReturn := make(chan JobsReturn, 100)
	claimBadgesUsecase := NewClaimBadgesUsecase(self.repository)
	claimRewardsUsecase := NewClaimRewardsUsecase(self.repository, self.mailProvider)

	for i := 0; i < numberOfThreads; i++ {
		go EventsWorker(self.jobs, jobsReturn, claimBadgesUsecase, claimRewardsUsecase, i)
	}

	for jobReturn := range jobsReturn {
		if jobReturn.Err != nil {
			self.jobs <- jobReturn.Job
		}
	}
}
