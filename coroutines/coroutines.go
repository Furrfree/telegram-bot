package coroutines

import (
	"context"
	"fmt"
	"time"

	"github.com/furrfree/telegram-bot/logger"
	"github.com/furrfree/telegram-bot/model"
	"github.com/furrfree/telegram-bot/service"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

var NewUserAddedChan = make(chan struct{})

func AddCoroutines(bot *telego.Bot) {
	go newAdmissionUser(bot)

}

func waitUntilTwoDaysAfterJoined(user model.NewUser) {
	targetTime := user.DateJoined.Add(10 * time.Second)
	duration := time.Until(targetTime)
	logger.Log(fmt.Sprintf("Coroutines: Waiting until %s", targetTime.String()))
	if duration > 0 {
		time.Sleep(duration)
	}
}

func newAdmissionUser(bot *telego.Bot) {
	logger.Log("Coroutines: Started newAdmissionUser")
	for {
		<-NewUserAddedChan // Wait until new user added

		// Get older user in DB and wait 2 days from its join
		olderUser := service.GetOlderJoinedNewUser()
		if olderUser != nil {
			waitUntilTwoDaysAfterJoined(*olderUser)
		}

		// If user still in admission group, ban it
		if service.UserExists(int64(olderUser.UserId)) {
			service.DeleteNewUser(olderUser.UserId)
			err := bot.BanChatMember(context.Background(), &telego.BanChatMemberParams{
				ChatID: tu.ID(int64(olderUser.ChatID)),
				UserID: int64(olderUser.UserId),
			})
			if err != nil {
				logger.Error(err)
				return
			}
			logger.Log(fmt.Sprintf("Removed user %s from Admission", olderUser.Username))
		}

	}
}
