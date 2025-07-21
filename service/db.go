package service

import (
	"errors"
	"sync"

	"github.com/furrfree/telegram-bot/logger"
	"github.com/furrfree/telegram-bot/model"
	"github.com/glebarez/sqlite"
	"github.com/lib/pq"
	"gorm.io/gorm"

	"time"
)

var lock = &sync.Mutex{}

var singleInstance *gorm.DB

func InitializeDb() {
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

	singleInstance.AutoMigrate(&model.Birthday{})
	singleInstance.AutoMigrate(&model.NewUser{})
}

func InsertBirthday(userId int64, groupId int64, birthday time.Time, username string) {
	singleInstance.Create(&model.Birthday{
		UserId:   int(userId),
		GroupId:  int(groupId),
		Date:     birthday,
		Username: username,
	})
}

func InsertNewUser(userId int64, username string, welcomeMessageId int) {
	singleInstance.Create(&model.NewUser{
		UserId:   int(userId),
		Username: username,
		Messages: pq.Int64Array{},
	})
}

func InsertNewUserMessage(userId int64, messageId int64) {
	var newUser model.NewUser
	singleInstance.Find(&newUser, "user_id=?", int(userId))
	newUser.Messages = append(newUser.Messages, messageId)
	singleInstance.Save(newUser)

}

func GetNearestBirthday(chatId int64) (*model.Birthday, error) {
	var nextBirthday model.Birthday

	var count int64
	singleInstance.Table("birthdays").Where("group_id", int(chatId)).Count(&count)

	if count == 0 {
		return nil, errors.New("No birthdays for this group")
	}

	today := time.Now().Format("01-02") // Format as MM-DD
	singleInstance.Raw("SELECT * FROM birthdays WHERE group_id = ? ORDER BY strftime('%m-%d',date) >= strftime('%m-%d',datetime('now') ) DESC, strftime('%m-%d',date ) ASC LIMIT 1", int(chatId), today).Scan(&nextBirthday)

	return &nextBirthday, nil

}

func GetNewUserFromUserId(userId int64) model.NewUser {

	var result model.NewUser

	singleInstance.Find(&result, "user_id=?", int(userId))

	return result

}

func GetNewUserByUsername(username string) model.NewUser {
	var result model.NewUser
	singleInstance.Where("username = ?", username).Find(&result)
	logger.Log(result)
	return result

}

func DeleteNewUser(newUserId int) {
	singleInstance.Where("user_id=?", newUserId).Delete(&model.NewUser{})
}
