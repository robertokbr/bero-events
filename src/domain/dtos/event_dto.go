package dtos

type event string

const (
	LOOT_COLLECTED   event = "loot_collected"
	REWARD_COLLECTED event = "reward_collected"
)

type EventDTO struct {
	UserID int64 `json:"userId"`
	Event  event `json:"event"`
}
