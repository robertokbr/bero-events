package dtos

type event string

const (
	LOOT_COLLECTED         event = "loot_collected"
	DAILY_INTERACTION_MADE event = "daily_interaction_made"
)

type CheckAchievementsDTO struct {
	UserID int64 `json:"userId"`
	Event  event `json:"event"`
}
