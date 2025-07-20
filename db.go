package main

import (
	"errors"

	"github.com/furrfree/telegram-bot/entities"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"time"
)

func setupDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("./database.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&entities.Birthday{})
	db.AutoMigrate(&entities.NewUser{})
	return db
}

func insertBirthday(db *gorm.DB, userId int64, groupId int64, birthday time.Time, username string) {

	db.Create(&entities.Birthday{
		UserId:   int(userId),
		GroupId:  int(groupId),
		Date:     birthday,
		Username: username,
	})
}

func insertNewUser(db *gorm.DB, userId int64, welcomeMessageId int) {
	db.Create(&entities.NewUser{
		UserId:           int(userId),
		WelcomeMessageId: welcomeMessageId,
	})
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

func getWelcomeMessageId(db *gorm.DB, userId int64) int {

	var result int

	db.Table("new_users").Where("user_id=?", int(userId)).Select("welcome_message_id").Scan(&result)

	return result

}

func deleteNewUser(db *gorm.DB, welcomeMessageId int) {
	db.Where("welcome_message_id=?", welcomeMessageId).Delete(&entities.NewUser{})
}
