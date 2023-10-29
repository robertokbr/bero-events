package repositories

import (
	"github.com/robertokbr/bero-events/src/domain/models"
	"github.com/robertokbr/bero-events/src/infra/database"
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
