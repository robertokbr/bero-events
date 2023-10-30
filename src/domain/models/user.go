package models

type User struct {
	Base
	ID    int64  `gorm:"column:id"`
	Email string `gorm:"column:email"`
	Name  string `gorm:"column:name"`
}
