package usecases

import (
	"fmt"
	"sync"

	"github.com/robertokbr/bero-events/src/domain/interfaces"
	"github.com/robertokbr/bero-events/src/domain/models"
	"github.com/robertokbr/bero-events/src/infra/database/repositories"
	"github.com/robertokbr/bero-events/src/logger"
)

type ClaimRewardsUsecase struct {
	repository   *repositories.MySqlRepository
	mailProvider interfaces.MailProvider
	mutex        *sync.Mutex
}

func NewClaimRewardsUsecase(repository *repositories.MySqlRepository, mailProvider interfaces.MailProvider) *ClaimRewardsUsecase {
	return &ClaimRewardsUsecase{
		repository:   repository,
		mailProvider: mailProvider,
		mutex:        &sync.Mutex{},
	}
}

func (self *ClaimRewardsUsecase) Execute(userID int64) error {
	self.mutex.Lock()
	user, err := self.repository.GetUserByID(userID)
	self.mutex.Unlock()

	if err != nil {
		logger.Errorf("Error while getting user %d: %s", userID, err.Error())
		return err
	}

	self.mutex.Lock()
	rewardIDs, err := self.repository.GetUserRewardIDsNotClaimedByUserID(userID)
	self.mutex.Unlock()

	if err != nil {
		logger.Errorf("Error while getting amount of achievement claims for user %d: %s", userID, err.Error())
		return err
	}

	if len(rewardIDs) == 0 {
		return nil
	}

	for _, rewardID := range rewardIDs {
		self.mutex.Lock()
		giftCard, err := self.repository.GetGiftCardByRewardID(rewardID)
		self.mutex.Unlock()

		if err != nil {
			logger.Errorf("Error while getting gift card for reward %d: %s", rewardID, err.Error())
			return err
		}

		if giftCard.ID == 0 {
			return nil
		}

		self.mutex.Lock()
		giftCard.IsUsed = true
		self.repository.UpdateGiftCard(&giftCard)
		userRewardPurchase := models.UserRewardPurchase{UserID: userID, RewardID: rewardID, IsDelivered: true}
		err = self.repository.UpdateUserRewardPurchase(&userRewardPurchase)
		self.mutex.Unlock()

		if err != nil {
			logger.Errorf("Error while updating user reward purchase %d: %s", userRewardPurchase.ID, err.Error())
			return err
		}

		emailBody := `
			<h1>Opa, %s!</h1>
			<img src="https://bero-bucket.s3.amazonaws.com/bero-money.png" />
			<p>Passando pra te entregar seu gift card! Espero que aproveite bastante :)</p>
			<p>Não esqueça de compartilhar no server pra todo mundo ver que você arrasa o/</p>
			<strong>#Código: %s</strong>
		`

		self.mutex.Lock()
		formattedEmailBody := fmt.Sprintf(emailBody, user.Name, giftCard.Code)
		err = self.mailProvider.Send(user.Email, "Plataforma do Bero - Gift Card", formattedEmailBody)
		self.mutex.Unlock()

		if err != nil {
			logger.Errorf("Error while sending email to user %d: %s", user.ID, err.Error())
			return err
		}

		logger.Debugf("Gift card %d delivered to user %d", giftCard.ID, user.ID)
	}

	return nil
}
