package badge_usecases

import (
	"github.com/robertokbr/bero-events/src/domain/enums"
	"github.com/robertokbr/bero-events/src/infra/database/repositories"
	"github.com/robertokbr/bero-events/src/logger"
)

type ClaimFirstHundredOfUsersAchievement struct {
	repository *repositories.MySqlRepository
}

func NewClaimFirstHundredOfUsersAchievement(repository *repositories.MySqlRepository) *ClaimFirstHundredOfUsersAchievement {
	return &ClaimFirstHundredOfUsersAchievement{
		repository: repository,
	}
}

func (self *ClaimFirstHundredOfUsersAchievement) Execute(userID int64) error {
	amountOfAchievementClaims, err := self.repository.GetAmountOfAchievementClaimsByAchievementID(int64(enums.FIRST_HUNDRED_OF_USERS))

	if err != nil {
		logger.Errorf("Error while getting amount of achievement claims for achievement %d: %s", enums.FIRST_HUNDRED_OF_USERS, err.Error())
		return err
	}

	if amountOfAchievementClaims > 101 {
		return nil
	}

	userAchievements, err := self.repository.GetUserAchievementsByUserAndAchievementID(
		userID,
		int64(enums.FIRST_HUNDRED_OF_USERS),
	)

	if err != nil {
		logger.Errorf("Error while getting user %d achievement %d: %s", userID, enums.FIRST_HUNDRED_OF_USERS, err.Error())
		return err
	}

	if len(userAchievements) > 0 {
		logger.Debugf("User %d has already collected achievement %d", userID, enums.FIRST_HUNDRED_OF_USERS)
		return nil
	}

	err = self.repository.CreateUserAchievement(userID, int64(enums.FIRST_HUNDRED_OF_USERS))

	if err != nil {
		logger.Errorf("Error while creating user %d achievement %d: %s", userID, enums.FIRST_HUNDRED_OF_USERS, err.Error())
		return err
	}

	logger.Debugf("Achievement %d unlocked for user %d", enums.FIRST_HUNDRED_OF_USERS, userID)

	return nil
}
