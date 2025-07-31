package service

import (
	"time"

	"github.com/furrfree/telegram-bot/coroutines"
	"github.com/furrfree/telegram-bot/database"
	"github.com/furrfree/telegram-bot/model"
	"github.com/lib/pq"
)

func GetOlderJoinedNewUser() *model.NewUser {
	var user model.NewUser
	database.Database.Order("date_joined asc").First(&user)
	return &user
}

func InsertNewUser(userId int64, chatId int64, username string) {
	database.Database.Create(&model.NewUser{
		UserId:     int(userId),
		Username:   username,
		Messages:   pq.Int64Array{},
		DateJoined: time.Now(),
		ChatID:     int(chatId),
	})
	// Send signal
	select {
	case coroutines.NewUserAddedChan <- struct{}{}:
	default:
	}
}

func InsertNewUserMessage(userId int64, messageId int64) {
	var newUser model.NewUser
	database.Database.Find(&newUser, "user_id=?", int(userId))
	newUser.Messages = append(newUser.Messages, messageId)
	database.Database.Save(newUser)

}

func UserExists(userId int64) bool {
	var user model.NewUser
	result := database.Database.First(&user, "user_id = ?", userId)
	return result.Error == nil
}

func GetNewUserFromUserId(userId int64) *model.NewUser {

	var result model.NewUser

	database.Database.Find(&result, "user_id=?", int(userId))

	return &result

}

func GetNewUserByUsername(username string) model.NewUser {
	var result model.NewUser
	database.Database.Where("username = ?", username).Find(&result)
	return result

}

func DeleteNewUser(newUserId int) {
	database.Database.Where("user_id=?", newUserId).Delete(&model.NewUser{})
}
