package main

import (
	"context"
	"fmt"
	"time"

	"github.com/furrfree/telegram-bot/commands"
	"github.com/furrfree/telegram-bot/configuration"
	"github.com/furrfree/telegram-bot/coroutines"
	"github.com/furrfree/telegram-bot/database"
	"github.com/furrfree/telegram-bot/handlers"
	"github.com/furrfree/telegram-bot/logger"
	"github.com/furrfree/telegram-bot/model"
	"github.com/furrfree/telegram-bot/service"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func WaitUntilTwoDaysAfterJoined(user model.NewUser) {
	targetTime := user.DateJoined.Add(48 * time.Hour)
	duration := time.Until(targetTime)
	fmt.Println(duration)
	if duration > 0 {
		time.Sleep(duration)
	}
}
func main() {
	database.InitializeDb()
	configuration.InitializeConfig()

	bot, botErr := telego.NewBot(configuration.Conf.Token)

	if botErr != nil {
		logger.Fatal(botErr)
	}

	// Get updates channel
	updates, _ := bot.UpdatesViaLongPolling(context.Background(), nil)

	// Create bot handler and specify from where to get updates
	bh, _ := th.NewBotHandler(bot, updates)

	// Stop handling updates
	defer func() { _ = bh.Stop() }()

	commands.AddCommands(bh, bot)
	handlers.AddHandlers(bh, bot)

	logger.Log("Bot started")

	go func() {
		for {
			<-coroutines.NewUserAddedChan // Wait until new user added

			// Get older user in DB and wait 2 days from its join
			olderUser := service.GetOlderJoinedNewUser()
			if olderUser != nil {
				WaitUntilTwoDaysAfterJoined(*olderUser)
			}

			fmt.Println("Finished waiting")

			// If user still in admission group, ban it
			if service.UserExists(int64(olderUser.UserId)) {
				fmt.Println("User exists")
				service.DeleteNewUser(olderUser.UserId)
				err := bot.BanChatMember(context.Background(), &telego.BanChatMemberParams{
					ChatID: tu.ID(int64(olderUser.ChatID)),
					UserID: int64(olderUser.UserId),
				})
				if err != nil {
					fmt.Println(err)
				}
			}

		}
	}()
	bh.Start()

}
