package tests

import (
	"testing"

	"github.com/furrfree/telegram-bot/internal/database"
	"github.com/furrfree/telegram-bot/internal/service"
)

func TestInsertNewUser(t *testing.T) {
	database.DropDatabase()
	database.InitializeDb()

	service.InsertNewUser(123, 456, "user1")

	user1 := service.GetNewUserFromUserId(123)
	if user1.Username != "user1" {
		t.Error(`Insert new_user failed `)
	}
}
