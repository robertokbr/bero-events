package badge_usecases

import (
	"sync"

	"github.com/robertokbr/bero-events/src/domain/enums"
	"github.com/robertokbr/bero-events/src/infra/database/repositories"
	"github.com/robertokbr/bero-events/src/logger"
)

type ClaimPopularAchievement struct {
	repository *repositories.MySqlRepository
	mutex      *sync.Mutex
}

func NewClaimPopularAchievement(repository *repositories.MySqlRepository) *ClaimPopularAchievement {
	return &ClaimPopularAchievement{
		repository: repository,
		mutex:      &sync.Mutex{},
	}
}

func (self *ClaimPopularAchievement) Execute(userID int64) error {
	self.mutex.Lock()
	amountOfGiftClaims, err := self.repository.GetUserAmountOfGiftClaimsByUserID(userID)
	self.mutex.Unlock()

	if err != nil {
		logger.Errorf("Error while getting amount of gift claims for user %d: %s", userID, err.Error())
		return err
	}

	if amountOfGiftClaims != 5 {
		return nil
	}

	self.mutex.Lock()
	userAchievements, err := self.repository.GetUserAchievementsByUserAndAchievementID(userID, int64(enums.USER_GIFT_CLAIMED_5_TIMES))
	self.mutex.Unlock()

	if err != nil {
		logger.Errorf("Error while getting user %d achievement %d: %s", userID, enums.USER_GIFT_CLAIMED_5_TIMES, err.Error())
		return err
	}

	if len(userAchievements) > 0 {
		logger.Debugf("User %d has already collected achievement %d", userID, enums.USER_GIFT_CLAIMED_5_TIMES)
		return nil
	}

	self.mutex.Lock()
	err = self.repository.CreateUserAchievement(userID, int64(enums.USER_GIFT_CLAIMED_5_TIMES))
	self.mutex.Unlock()

	if err != nil {
		logger.Errorf("Error while creating user %d achievement %d: %s", userID, enums.USER_GIFT_CLAIMED_5_TIMES, err.Error())
		return err
	}

	logger.Debugf("Achievement %d unlocked for user %d", enums.USER_GIFT_CLAIMED_5_TIMES, userID)

	return nil
}
