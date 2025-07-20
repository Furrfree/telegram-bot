package main

import (
	"errors"
	"sync"

	"github.com/furrfree/telegram-bot/entities"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"time"
)

var lock = &sync.Mutex{}

var singleInstance *gorm.DB

func initializeDb() {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			setupDb()
		}
	}
}

func setupDb() {
	var err error
	singleInstance, err = gorm.Open(sqlite.Open("./database.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	singleInstance.AutoMigrate(&entities.Birthday{})
	singleInstance.AutoMigrate(&entities.NewUser{})
}

func insertBirthday(userId int64, groupId int64, birthday time.Time, username string) {
	singleInstance.Create(&entities.Birthday{
		UserId:   int(userId),
		GroupId:  int(groupId),
		Date:     birthday,
		Username: username,
	})
}

func insertNewUser(userId int64, welcomeMessageId int) {
	singleInstance.Create(&entities.NewUser{
		UserId:           int(userId),
		WelcomeMessageId: welcomeMessageId,
	})
}

func getNearestBirthday(chatId int64) (*entities.Birthday, error) {
	var nextBirthday entities.Birthday

	var count int64
	singleInstance.Table("birthdays").Where("group_id", int(chatId)).Count(&count)

	if count == 0 {
		return nil, errors.New("No birthdays for this group")
	}

	today := time.Now().Format("01-02") // Format as MM-DD
	singleInstance.Raw("SELECT * FROM birthdays WHERE group_id = ? ORDER BY strftime('%m-%d',date) >= strftime('%m-%d',datetime('now') ) DESC, strftime('%m-%d',date ) ASC LIMIT 1", int(chatId), today).Scan(&nextBirthday)

	return &nextBirthday, nil

}

func getWelcomeMessageId(userId int64) int {

	var result int

	singleInstance.Table("new_users").Where("user_id=?", int(userId)).Select("welcome_message_id").Scan(&result)

	return result

}

func deleteNewUser(welcomeMessageId int) {
	singleInstance.Where("welcome_message_id=?", welcomeMessageId).Delete(&entities.NewUser{})
}
