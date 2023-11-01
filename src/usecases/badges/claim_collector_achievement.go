package badge_usecases

import (
	"sync"

	"github.com/robertokbr/events-worker/src/domain/enums"
	"github.com/robertokbr/events-worker/src/infra/database/repositories"
	"github.com/robertokbr/events-worker/src/logger"
)

type ClaimCollectorAchievement struct {
	repository *repositories.MySqlRepository
	mutex      *sync.Mutex
}

func NewClaimCollectorAchievement(repository *repositories.MySqlRepository) *ClaimCollectorAchievement {
	return &ClaimCollectorAchievement{
		repository: repository,
		mutex:      &sync.Mutex{},
	}
}

func (self *ClaimCollectorAchievement) Execute(userID int64) error {
	self.mutex.Lock()
	amountOfPureLoots, err := self.repository.GetUserAmountOfPureLootsByUserID(userID)
	self.mutex.Unlock()

	if err != nil {
		logger.Errorf("Error while getting amount of pure loots for user %d: %s", userID, err.Error())
		return err
	}

	if amountOfPureLoots == 10 {
		self.mutex.Lock()
		userAchievements, err := self.repository.GetUserAchievementsByUserAndAchievementID(userID, int64(enums.USER_COLLECTED_10_PURE_LOOTS))
		self.mutex.Unlock()

		if err != nil {
			logger.Errorf("Error while getting user %d achievement %d: %s", userID, enums.USER_COLLECTED_10_PURE_LOOTS, err.Error())
			return err
		}

		if len(userAchievements) > 0 {
			logger.Debugf("User %d has already collected achievement %d", userID, enums.USER_COLLECTED_10_PURE_LOOTS)
			return nil
		}

		self.mutex.Lock()
		err = self.repository.CreateUserAchievement(userID, int64(enums.USER_COLLECTED_10_PURE_LOOTS))
		self.mutex.Unlock()

		if err != nil {
			logger.Errorf("Error while creating user %d achievement %d: %s", userID, enums.USER_COLLECTED_10_PURE_LOOTS, err.Error())
			return err
		}

		logger.Debugf("Achievement %d unlocked for user %d", enums.USER_COLLECTED_10_PURE_LOOTS, userID)
	}

	return nil
}
