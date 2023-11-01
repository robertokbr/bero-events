package repositories

import (
	"github.com/robertokbr/events-worker/src/domain/models"
	"github.com/robertokbr/events-worker/src/infra/database"
)

type MySqlRepository struct {
	db *database.MysqlDB
}

func NewMySqlRepository(connection *database.MysqlDB) *MySqlRepository {
	return &MySqlRepository{
		db: connection,
	}
}

func (self *MySqlRepository) GetUserAmountOfPureLootsByUserID(userID int64) (int64, error) {
	userLoots := make(map[string]interface{})

	tx := self.db.DB.Raw(
		`
			SELECT count(*) as collected_loots FROM UsersLoots ul
			JOIN Loots l on l.id = ul.lootId
			WHERE ul.userId = ?
			AND l.isGift = false
		`,
		userID,
	).Scan(&userLoots)

	if tx.Error != nil {
		if tx.Error.Error() == "record not found" {
			return -1, nil
		} else {
			return 0, tx.Error
		}
	}

	return userLoots["collected_loots"].(int64), nil
}

func (self *MySqlRepository) GetUserAchievementsByUserAndAchievementID(userID, achievementID int64) ([]models.Achievement, error) {
	query := `
		SELECT * FROM UsersAchievements 
		WHERE userId = ? 
		AND achievementId = ?
	`

	userAchievements := make([]models.Achievement, 0)

	tx := self.db.DB.Raw(
		query,
		userID,
		achievementID,
	).Scan(&userAchievements)

	if tx.Error != nil {
		if tx.Error.Error() == "record not found" {
			return nil, nil
		}

		return nil, tx.Error
	}

	return userAchievements, nil
}

func (self *MySqlRepository) GetUserAmountOfPureGemsByUserID(userID int64) (int64, error) {
	userGems := make(map[string]interface{})

	tx := self.db.DB.Raw(
		`
			SELECT count(*) as gems FROM Loots l 
			WHERE l.isGift = 0
			AND l.firstClaimedBy = ?
		`,
		userID,
	).Scan(&userGems)

	if tx.Error != nil {
		if tx.Error.Error() == "record not found" {
			return -1, nil
		} else {
			return 0, tx.Error
		}
	}

	return userGems["gems"].(int64), nil
}

func (self *MySqlRepository) GetUserAmountOfGiftClaimsByUserID(userID int64) (int64, error) {
	amountOfClaims := make(map[string]interface{})

	tx := self.db.DB.Raw(
		`
			SELECT count(*) AS amount_of_claims FROM UsersLoots ul
			JOIN UserGiftLoots ugl on ugl.lootId = ul.lootId 
			WHERE ugl.userId = ?
		`,
		userID,
	).Scan(&amountOfClaims)

	if tx.Error != nil {
		if tx.Error.Error() == "record not found" {
			return -1, nil
		} else {
			return 0, tx.Error
		}
	}

	return amountOfClaims["amount_of_claims"].(int64), nil
}

func (self *MySqlRepository) GetAmountOfAchievementClaimsByAchievementID(achievementID int64) (int64, error) {
	amountOfAchievementClaims := make(map[string]interface{})

	tx := self.db.DB.Raw(
		`
			SELECT count(ua.userId) claims FROM UsersAchievements ua 
			WHERE ua.achievementId = ?
			GROUP by ua.userId 
			HAVING claims = 1
		`,
		achievementID,
	).Scan(&amountOfAchievementClaims)

	if tx.Error != nil {
		if tx.Error.Error() == "record not found" {
			return -1, nil
		} else {
			return 0, tx.Error
		}
	}

	return amountOfAchievementClaims["claims"].(int64), nil
}

func (self *MySqlRepository) CreateUserAchievement(userID, achievementID int64) error {
	tx := self.db.DB.Table("UsersAchievements").Create(&models.UserAchievement{
		UserID:        userID,
		AchievementID: achievementID,
	})

	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (self *MySqlRepository) GetUserRewardIDsNotClaimedByUserID(userID int64) ([]int64, error) {
	IDs := make([]int64, 0)

	const query = `
		SELECT urp.rewardId FROM UsersRewardPurchases urp 
		WHERE urp.userId = ?
		AND urp.isDelivered = false
	`
	tx := self.db.DB.Raw(query, userID).Pluck("rewardId", &IDs)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return IDs, nil
}

func (self *MySqlRepository) GetGiftCardByRewardID(giftCardID int64) (models.GiftCard, error) {
	giftCard := models.GiftCard{}

	const query = `
		SELECT gc.code, gc.id FROM rewards_gift_cards rgc 
		JOIN gift_cards gc ON gc.id = rgc.id 
		WHERE rgc.reward_id = ?
		AND gc.is_used = FALSE
	`

	tx := self.db.DB.Raw(query, giftCardID).Scan(&giftCard)

	if tx.Error != nil {
		return models.GiftCard{}, tx.Error
	}

	return giftCard, nil
}

func (self *MySqlRepository) UpdateGiftCard(giftCard *models.GiftCard) error {
	tx := self.db.DB.Table("gift_cards").Where("id = ?", giftCard.ID).Updates(&giftCard)

	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (self *MySqlRepository) UpdateUserRewardPurchase(userRewardPurchase *models.UserRewardPurchase) error {
	tx := self.db.DB.Table("UsersRewardPurchases").Where(
		"rewardId = ? AND userId = ?",
		userRewardPurchase.RewardID,
		userRewardPurchase.UserID,
	).Updates(&userRewardPurchase)

	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (self *MySqlRepository) GetUserByID(userID int64) (models.User, error) {
	user := models.User{}

	tx := self.db.DB.Table("Users").Where("id = ?", userID).Scan(&user)

	if tx.Error != nil {
		return models.User{}, tx.Error
	}

	return user, nil
}
