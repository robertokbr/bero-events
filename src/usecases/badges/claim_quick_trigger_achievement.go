package badge_usecases

import (
	"sync"

	"github.com/robertokbr/events-worker/src/domain/enums"
	"github.com/robertokbr/events-worker/src/infra/database/repositories"
	"github.com/robertokbr/events-worker/src/logger"
)

type ClaimQuickTriggerAchievement struct {
	repository *repositories.MySqlRepository
	mutex      *sync.Mutex
}

func NewClaimQuickTriggerAchievement(repository *repositories.MySqlRepository) *ClaimQuickTriggerAchievement {
	return &ClaimQuickTriggerAchievement{
		repository: repository,
		mutex:      &sync.Mutex{},
	}
}

func (self *ClaimQuickTriggerAchievement) Execute(userID int64) error {
	self.mutex.Lock()
	amountOfPureGems, err := self.repository.GetUserAmountOfPureGemsByUserID(userID)
	self.mutex.Unlock()

	if err != nil {
		logger.Errorf("Error while getting amount of pure gems for user %d: %s", userID, err.Error())
		return err
	}

	if amountOfPureGems == 1 {
		self.mutex.Lock()
		userAchievements, err := self.repository.GetUserAchievementsByUserAndAchievementID(userID, int64(enums.FIRST_COLLECTED_PURE_GEM))
		self.mutex.Unlock()

		if err != nil {
			logger.Errorf("Error while getting user %d achievement %d: %s", userID, enums.FIRST_COLLECTED_PURE_GEM, err.Error())
			return err
		}

		if len(userAchievements) > 0 {
			logger.Debugf("User %d has already collected achievement %d", userID, enums.FIRST_COLLECTED_PURE_GEM)
			return nil
		}

		self.mutex.Lock()
		err = self.repository.CreateUserAchievement(userID, int64(enums.FIRST_COLLECTED_PURE_GEM))
		self.mutex.Unlock()

		if err != nil {
			logger.Errorf("Error while creating user %d achievement %d: %s", userID, enums.FIRST_COLLECTED_PURE_GEM, err.Error())
			return err
		}

		logger.Debugf("Achievement %d unlocked for user %d", enums.FIRST_COLLECTED_PURE_GEM, userID)
	}

	return nil
}
