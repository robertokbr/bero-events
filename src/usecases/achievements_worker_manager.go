package usecases

import (
	"github.com/robertokbr/bero-events/src/domain/dtos"
	"github.com/robertokbr/bero-events/src/infra/database/repositories"
)

type JobsReturn struct {
	Job *dtos.CheckAchievementsDTO
	Err error
}

type AchievementsWorkerManager struct {
	repository *repositories.MySqlRepository
	jobs       chan *dtos.CheckAchievementsDTO
}

func NewAchievementsWorkerManager(jobs chan *dtos.CheckAchievementsDTO, repository *repositories.MySqlRepository) *AchievementsWorkerManager {
	return &AchievementsWorkerManager{
		jobs:       jobs,
		repository: repository,
	}
}

func (self *AchievementsWorkerManager) Start(numberOfThreads int) {
	jobsReturn := make(chan JobsReturn, 100)
	badgesUsecase := NewBadgesUsecase(self.repository)

	for i := 0; i < numberOfThreads; i++ {
		go AchievementsWorker(self.jobs, jobsReturn, badgesUsecase, i)
	}

	for jobReturn := range jobsReturn {
		if jobReturn.Err != nil {
			self.jobs <- jobReturn.Job
		}
	}
}
