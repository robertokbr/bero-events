package database

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/robertokbr/bero-events/src/domain/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MysqlDB struct {
	DB *gorm.DB
}

func NewFakeDB() *MysqlDB {
	db, err := gorm.Open(sqlite.Open(":memory"), &gorm.Config{})

	if err != nil {
		panic("Database connection failed")
	}

	db.AutoMigrate(
		&models.Achievement{},
		&models.Loot{},
		&models.UserAchievement{},
		&models.UserGiftLoot{},
		&models.UserLoot{},
	)

	return &MysqlDB{
		DB: db,
	}
}

func NewMysqlDB() *MysqlDB {
	godotenv.Load()

	dns := os.Getenv("MYSQL_DNS")

	if dns == "" {
		panic("MYSQL_DNS environment variable not set")
	}

	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})

	if err != nil {
		panic("Database connection failed")
	}

	return &MysqlDB{
		DB: db,
	}
}

func (self *MysqlDB) Close() error {
	db, err := self.DB.DB()

	if err != nil {
		return err
	}

	return db.Close()
}

var MySQL = NewMysqlDB()
