package tests

import (
	"testing"

	"github.com/furrfree/telegram-bot/database"
	"github.com/furrfree/telegram-bot/service"
)

func TestInsertNewUser(t *testing.T) {
	database.DropDatabase()
	database.InitializeDb()

	service.InsertNewUser(123, 456, "user1")

	user1 := service.GetNewUserByUsername("user1")
	if user1.Username != "user1" {
		t.Error(`Insert new_user failed `)
	}
}
