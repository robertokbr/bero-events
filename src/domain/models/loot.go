package models

type Loot struct {
	Base
	Gems   string `json:"gems"`
	IsGift bool   `json:"isGift"`
}
