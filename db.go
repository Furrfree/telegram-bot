package main

import (
	"errors"

	"github.com/furrfree/telegram-bot/entities"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"time"
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

	db.Create(&entities.Birthday{
		UserId:   birthday.UserId,
		GroupId:  birthday.GroupId,
		Date:     birthday.Date,
		Username: birthday.Username,
	})
	db.Create(birthday)
}

func getNearestBirthday(db *gorm.DB, chatId int64) (*entities.Birthday, error) {
	var nextBirthday entities.Birthday

	var count int64
	db.Table("birthdays").Where("group_id", int(chatId)).Count(&count)

	if count == 0 {
		return nil, errors.New("No birthdays for this group")
	}

	today := time.Now().Format("01-02") // Format as MM-DD
	db.Raw("SELECT * FROM birthdays WHERE group_id = ? ORDER BY strftime('%m-%d',date) >= strftime('%m-%d',datetime('now') ) DESC, strftime('%m-%d',date ) ASC LIMIT 1", int(chatId), today).Scan(&nextBirthday)

	return &nextBirthday, nil

}
