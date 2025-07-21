package service

import (
	"errors"
	"time"

	"github.com/furrfree/telegram-bot/database"
	"github.com/furrfree/telegram-bot/model"
)

func InsertBirthday(userId int64, groupId int64, birthday time.Time, username string) {
	database.Database.Create(&model.Birthday{
		UserId:   int(userId),
		GroupId:  int(groupId),
		Date:     birthday,
		Username: username,
	})
}

func GetNearestBirthday(chatId int64) (*model.Birthday, error) {
	var nextBirthday model.Birthday

	var count int64
	database.Database.Table("birthdays").Where("group_id", int(chatId)).Count(&count)

	if count == 0 {
		return nil, errors.New("No birthdays for this group")
	}

	today := time.Now().Format("01-02") // Format as MM-DD
	database.Database.Raw("SELECT * FROM birthdays WHERE group_id = ? ORDER BY strftime('%m-%d',date) >= strftime('%m-%d',datetime('now') ) DESC, strftime('%m-%d',date ) ASC LIMIT 1", int(chatId), today).Scan(&nextBirthday)

	return &nextBirthday, nil

}
