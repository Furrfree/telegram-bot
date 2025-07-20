package main

import (
	"github.com/furrfree/telegram-bot/entities"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("./test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&entities.Birthday{})
	return db
}

func insertBirthday(db *gorm.DB, birthday *entities.Birthday) {
	db.Create(birthday)
}
