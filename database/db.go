package database

import (
	"sync"

	"github.com/furrfree/telegram-bot/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var lock = &sync.Mutex{}

var Database *gorm.DB

func InitializeDb() {
	if Database == nil {
		lock.Lock()
		defer lock.Unlock()
		if Database == nil {
			setupDb()
		}
	}
}

func setupDb() {
	var err error
	Database, err = gorm.Open(sqlite.Open("./database.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	Database.AutoMigrate(&model.Birthday{})
	Database.AutoMigrate(&model.NewUser{})
}
