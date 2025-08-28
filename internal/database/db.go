package database

import (
	"fmt"
	"os"
	"sync"

	"github.com/furrfree/telegram-bot/internal/logger"
	"github.com/furrfree/telegram-bot/internal/model"
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

// DropDatabase closes the database connection and deletes the database file.
// It safely handles cases where the database might not be initialized.
func DropDatabase() {
	lock.Lock()
	defer lock.Unlock()

	if Database != nil {
		sqlDB, err := Database.DB()
		if err == nil {
			sqlDB.Close()
		}
		Database = nil
	}

	// Remove the database file
	err := os.Remove("./database.db")
	if err != nil && !os.IsNotExist(err) {
		logger.Error(fmt.Sprintf("Error removing database: %d", err))
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
