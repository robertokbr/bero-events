package usecases

import (
	"github.com/robertokbr/bero-events/src/infra/database/repositories"
	badgeUsecases "github.com/robertokbr/bero-events/src/usecases/badges"
)

type ClaimBadgesUsecase struct {
	claimCollectorAchievement            *badgeUsecases.ClaimCollectorAchievement
	claimFirstHundredOfUsersAchievement  *badgeUsecases.ClaimFirstHundredOfUsersAchievement
	claimFirstThousandOfUsersAchievement *badgeUsecases.ClaimFirstThousandOfUsersAchievement
	claimInfluencerAchievement           *badgeUsecases.ClaimInfluencerAchievement
	claimPopularAchievement              *badgeUsecases.ClaimPopularAchievement
	claimQuickTriggerAchievement         *badgeUsecases.ClaimQuickTriggerAchievement
}

func NewClaimBadgesUsecase(repository *repositories.MySqlRepository) *ClaimBadgesUsecase {
	return &ClaimBadgesUsecase{
		claimCollectorAchievement:            badgeUsecases.NewClaimCollectorAchievement(repository),
		claimFirstHundredOfUsersAchievement:  badgeUsecases.NewClaimFirstHundredOfUsersAchievement(repository),
		claimFirstThousandOfUsersAchievement: badgeUsecases.NewClaimFirstThousandOfUsersAchievement(repository),
		claimInfluencerAchievement:           badgeUsecases.NewClaimInfluencerAchievement(repository),
		claimPopularAchievement:              badgeUsecases.NewClaimPopularAchievement(repository),
		claimQuickTriggerAchievement:         badgeUsecases.NewClaimQuickTriggerAchievement(repository),
	}
}

func (self *ClaimBadgesUsecase) Execute(userID int64) {
	self.claimCollectorAchievement.Execute(userID)
	self.claimFirstHundredOfUsersAchievement.Execute(userID)
	self.claimFirstThousandOfUsersAchievement.Execute(userID)
	self.claimInfluencerAchievement.Execute(userID)
	self.claimPopularAchievement.Execute(userID)
	self.claimQuickTriggerAchievement.Execute(userID)
}
