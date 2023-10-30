package models

type GiftCard struct {
	Base
	Code   string
	IsUsed bool
}
