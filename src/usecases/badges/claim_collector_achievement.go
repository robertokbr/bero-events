package badge_usecases

import (
	"github.com/robertokbr/bero-events/src/domain/enums"
	"github.com/robertokbr/bero-events/src/infra/database/repositories"
	"github.com/robertokbr/bero-events/src/logger"
)

type ClaimCollectorAchievement struct {
	repository *repositories.MySqlRepository
}

func NewClaimCollectorAchievement(repository *repositories.MySqlRepository) *ClaimCollectorAchievement {
	return &ClaimCollectorAchievement{
		repository: repository,
	}
}

func (self *ClaimCollectorAchievement) Execute(userID int64) error {
	amountOfPureLoots, err := self.repository.GetUserAmountOfPureLootsByUserID(userID)

	if err != nil {
		logger.Errorf("Error while getting amount of pure loots for user %d: %s", userID, err.Error())
		return err
	}

	if amountOfPureLoots == 10 {
		userAchievements, err := self.repository.GetUserAchievementsByUserAndAchievementID(
			userID,
			int64(enums.USER_COLLECTED_10_PURE_LOOTS),
		)

		if err != nil {
			logger.Errorf("Error while getting user %d achievement %d: %s", userID, enums.USER_COLLECTED_10_PURE_LOOTS, err.Error())
			return err
		}

		if len(userAchievements) > 0 {
			logger.Debugf("User %d has already collected achievement %d", userID, enums.USER_COLLECTED_10_PURE_LOOTS)
			return nil
		}

		err = self.repository.CreateUserAchievement(userID, int64(enums.USER_COLLECTED_10_PURE_LOOTS))

		if err != nil {
			logger.Errorf("Error while creating user %d achievement %d: %s", userID, enums.USER_COLLECTED_10_PURE_LOOTS, err.Error())
			return err
		}

		logger.Debugf("Achievement %d unlocked for user %d", enums.USER_COLLECTED_10_PURE_LOOTS, userID)
	}

	return nil
}
