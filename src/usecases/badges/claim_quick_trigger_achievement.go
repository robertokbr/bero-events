package badge_usecases

import (
	"github.com/robertokbr/bero-events/src/domain/enums"
	"github.com/robertokbr/bero-events/src/infra/database/repositories"
	"github.com/robertokbr/bero-events/src/logger"
)

type ClaimQuickTriggerAchievement struct {
	repository *repositories.MySqlRepository
}

func NewClaimQuickTriggerAchievement(repository *repositories.MySqlRepository) *ClaimQuickTriggerAchievement {
	return &ClaimQuickTriggerAchievement{
		repository: repository,
	}
}

func (self *ClaimQuickTriggerAchievement) Execute(userID int64) error {
	amountOfPureGems, err := self.repository.GetUserAmountOfPureGemsByUserID(userID)

	if err != nil {
		logger.Errorf("Error while getting amount of pure gems for user %d: %s", userID, err.Error())
		return err
	}

	if amountOfPureGems == 1 {
		userAchievements, err := self.repository.GetUserAchievementsByUserAndAchievementID(
			userID,
			int64(enums.FIRST_COLLECTED_PURE_GEM),
		)

		if err != nil {
			logger.Errorf("Error while getting user %d achievement %d: %s", userID, enums.FIRST_COLLECTED_PURE_GEM, err.Error())
			return err
		}

		if len(userAchievements) > 0 {
			logger.Debugf("User %d has already collected achievement %d", userID, enums.FIRST_COLLECTED_PURE_GEM)
			return nil
		}

		err = self.repository.CreateUserAchievement(userID, int64(enums.FIRST_COLLECTED_PURE_GEM))

		if err != nil {
			logger.Errorf("Error while creating user %d achievement %d: %s", userID, enums.FIRST_COLLECTED_PURE_GEM, err.Error())
			return err
		}

		logger.Debugf("Achievement %d unlocked for user %d", enums.FIRST_COLLECTED_PURE_GEM, userID)
	}

	return nil
}
