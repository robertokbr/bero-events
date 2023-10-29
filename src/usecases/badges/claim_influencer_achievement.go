package badge_usecases

import (
	"github.com/robertokbr/bero-events/src/domain/enums"
	"github.com/robertokbr/bero-events/src/infra/database/repositories"
	"github.com/robertokbr/bero-events/src/logger"
)

type ClaimInfluencerAchievement struct {
	repository *repositories.MySqlRepository
}

func NewClaimInfluencerAchievement(repository *repositories.MySqlRepository) *ClaimInfluencerAchievement {
	return &ClaimInfluencerAchievement{
		repository: repository,
	}
}

func (self *ClaimInfluencerAchievement) Execute(userID int64) error {
	amountOfGiftClaims, err := self.repository.GetUserAmountOfGiftClaimsByUserID(userID)

	if err != nil {
		logger.Errorf("Error while getting amount of gift claims for user %d: %s", userID, err.Error())
		return err
	}

	if amountOfGiftClaims != 10 {
		return nil
	}

	userAchievements, err := self.repository.GetUserAchievementsByUserAndAchievementID(
		userID,
		int64(enums.USER_GIFT_CLAIMED_10_TIMES),
	)

	if err != nil {
		logger.Errorf("Error while getting user %d achievement %d: %s", userID, enums.USER_GIFT_CLAIMED_10_TIMES, err.Error())
		return err
	}

	if len(userAchievements) > 0 {
		logger.Debugf("User %d has already collected achievement %d", userID, enums.USER_GIFT_CLAIMED_10_TIMES)
		return nil
	}

	err = self.repository.CreateUserAchievement(userID, int64(enums.USER_GIFT_CLAIMED_10_TIMES))

	if err != nil {
		logger.Errorf("Error while creating user %d achievement %d: %s", userID, enums.USER_GIFT_CLAIMED_10_TIMES, err.Error())
		return err
	}

	logger.Debugf("Achievement %d unlocked for user %d", enums.USER_GIFT_CLAIMED_10_TIMES, userID)

	return nil
}
