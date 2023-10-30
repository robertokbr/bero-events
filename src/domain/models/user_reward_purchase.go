package models

type UserRewardPurchase struct {
	Base
	UserID      int64 `gorm:"column:userId"`
	RewardID    int64 `gorm:"column:rewardId"`
	IsDelivered bool  `gorm:"column:isDelivered"`
}
