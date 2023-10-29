package models

type UserAchievement struct {
	Base
	UserID        int64 `json:"userId" gorm:"column:userId"`
	AchievementID int64 `json:"achievementId" gorm:"column:achievementId"`
}
