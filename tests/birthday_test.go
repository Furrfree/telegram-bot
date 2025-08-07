package tests

import (
	"testing"
	"time"

	"github.com/furrfree/telegram-bot/database"
	"github.com/furrfree/telegram-bot/service"
)

func TestInsertBirthday(t *testing.T) {
	database.DropDatabase()
	database.InitializeDb()
	service.InsertBirthday(123, 456, time.Date(2002, 3, 11, 0, 0, 0, 0, time.UTC), "username")
	nearestBirthday, err := service.GetNearestBirthday(456)
	if err != nil {
		t.Fatal(err)
	}
	if nearestBirthday.Username != "username" {
		t.Error(`Insert Birthday failed `)
	}
}

func TestGetNearestBirthday(t *testing.T) {
	// TODO
	//service.InsertBirthday(123, 456, time.Date(2005, 1, 11, 0, 0, 0, 0, time.UTC), "latestBi")

}
