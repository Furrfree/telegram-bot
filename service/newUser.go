package service

import (
	"github.com/furrfree/telegram-bot/database"
	"github.com/furrfree/telegram-bot/logger"
	"github.com/furrfree/telegram-bot/model"
	"github.com/lib/pq"
)

func InsertNewUser(userId int64, username string, welcomeMessageId int) {
	database.Database.Create(&model.NewUser{
		UserId:   int(userId),
		Username: username,
		Messages: pq.Int64Array{},
	})
}

func InsertNewUserMessage(userId int64, messageId int64) {
	var newUser model.NewUser
	database.Database.Find(&newUser, "user_id=?", int(userId))
	newUser.Messages = append(newUser.Messages, messageId)
	database.Database.Save(newUser)

}
func GetNewUserFromUserId(userId int64) model.NewUser {

	var result model.NewUser

	database.Database.Find(&result, "user_id=?", int(userId))

	return result

}

func GetNewUserByUsername(username string) model.NewUser {
	var result model.NewUser
	database.Database.Where("username = ?", username).Find(&result)
	logger.Log(result)
	return result

}

func DeleteNewUser(newUserId int) {
	database.Database.Where("user_id=?", newUserId).Delete(&model.NewUser{})
}
