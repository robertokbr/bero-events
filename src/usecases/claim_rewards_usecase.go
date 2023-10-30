package usecases

import (
	"fmt"

	"github.com/robertokbr/bero-events/src/domain/interfaces"
	"github.com/robertokbr/bero-events/src/domain/models"
	"github.com/robertokbr/bero-events/src/infra/database/repositories"
	"github.com/robertokbr/bero-events/src/logger"
)

type ClaimRewardsUsecase struct {
	repository   *repositories.MySqlRepository
	mailProvider interfaces.MailProvider
}

func NewClaimRewardsUsecase(repository *repositories.MySqlRepository, mailProvider interfaces.MailProvider) *ClaimRewardsUsecase {
	return &ClaimRewardsUsecase{
		repository:   repository,
		mailProvider: mailProvider,
	}
}

func (self *ClaimRewardsUsecase) Execute(userID int64) error {
	user, err := self.repository.GetUserByID(userID)

	if err != nil {
		logger.Errorf("Error while getting user %d: %s", userID, err.Error())
		return err
	}

	rewardIDs, err := self.repository.GetUserRewardIDsNotClaimedByUserID(userID)

	if err != nil {
		logger.Errorf("Error while getting amount of achievement claims for user %d: %s", userID, err.Error())
		return err
	}

	if len(rewardIDs) == 0 {
		return nil
	}

	for _, rewardID := range rewardIDs {
		giftCard, err := self.repository.GetGiftCardByRewardID(rewardID)

		if err != nil {
			logger.Errorf("Error while getting gift card for reward %d: %s", rewardID, err.Error())
			return err
		}

		if giftCard.ID == 0 {
			return nil
		}

		giftCard.IsUsed = true

		self.repository.UpdateGiftCard(&giftCard)

		userRewardPurchase := models.UserRewardPurchase{
			UserID:      userID,
			RewardID:    rewardID,
			IsDelivered: true,
		}

		err = self.repository.UpdateUserRewardPurchase(&userRewardPurchase)

		if err != nil {
			logger.Errorf("Error while updating user reward purchase %d: %s", userRewardPurchase.ID, err.Error())
			return err
		}

		emailBody := `
			<h1>Opa, %s!</h1>
			<img src="https://bero-bucket.s3.amazonaws.com/bero-money.png" />
			<p>Passando pra te entregar seu gift card! Espero que aproveite bastante :)</p>
			<p>Não esqueça de compartilhar no server pra todo mundo ver que você arrasa o/</p>
			<strong>- Código: %s</strong>
		`

		formattedEmailBody := fmt.Sprintf(emailBody, user.Name, giftCard.Code)

		err = self.mailProvider.Send(
			user.Email,
			"Plataforma do Bero - Gift Card",
			formattedEmailBody,
		)

		if err != nil {
			logger.Errorf("Error while sending email to user %d: %s", user.ID, err.Error())
			return err
		}

		logger.Debugf("Gift card %d delivered to user %d", giftCard.ID, user.ID)
	}

	return nil
}
