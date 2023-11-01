package badge_usecases

import (
	"sync"

	"github.com/robertokbr/events-worker/src/domain/enums"
	"github.com/robertokbr/events-worker/src/infra/database/repositories"
	"github.com/robertokbr/events-worker/src/logger"
)

type ClaimFirstHundredOfUsersAchievement struct {
	repository *repositories.MySqlRepository
	mutex      *sync.Mutex
}

func NewClaimFirstHundredOfUsersAchievement(repository *repositories.MySqlRepository) *ClaimFirstHundredOfUsersAchievement {
	return &ClaimFirstHundredOfUsersAchievement{
		repository: repository,
		mutex:      &sync.Mutex{},
	}
}

func (self *ClaimFirstHundredOfUsersAchievement) Execute(userID int64) error {
	self.mutex.Lock()
	amountOfPureLoots, err := self.repository.GetUserAmountOfPureLootsByUserID(userID)
	self.mutex.Unlock()

	if err != nil {
		logger.Errorf("Error while getting user %d amount of pure loots: %s", userID, err.Error())
		return err
	}

	if amountOfPureLoots > 1 {
		return nil
	}

	self.mutex.Lock()
	amountOfAchievementClaims, err := self.repository.GetAmountOfAchievementClaimsByAchievementID(int64(enums.FIRST_HUNDRED_OF_USERS))
	self.mutex.Unlock()

	if err != nil {
		logger.Errorf("Error while getting amount of achievement claims for achievement %d: %s", enums.FIRST_HUNDRED_OF_USERS, err.Error())
		return err
	}

	if amountOfAchievementClaims > 101 {
		return nil
	}

	self.mutex.Lock()
	userAchievements, err := self.repository.GetUserAchievementsByUserAndAchievementID(userID, int64(enums.FIRST_HUNDRED_OF_USERS))
	self.mutex.Unlock()

	if err != nil {
		logger.Errorf("Error while getting user %d achievement %d: %s", userID, enums.FIRST_HUNDRED_OF_USERS, err.Error())
		return err
	}

	if len(userAchievements) > 0 {
		logger.Debugf("User %d has already collected achievement %d", userID, enums.FIRST_HUNDRED_OF_USERS)
		return nil
	}

	self.mutex.Lock()
	err = self.repository.CreateUserAchievement(userID, int64(enums.FIRST_HUNDRED_OF_USERS))
	self.mutex.Unlock()

	if err != nil {
		logger.Errorf("Error while creating user %d achievement %d: %s", userID, enums.FIRST_HUNDRED_OF_USERS, err.Error())
		return err
	}

	logger.Debugf("Achievement %d unlocked for user %d", enums.FIRST_HUNDRED_OF_USERS, userID)

	return nil
}
