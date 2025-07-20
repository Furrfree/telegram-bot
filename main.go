package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func main() {
	// Get Bot token from environment variables

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	botToken := os.Getenv("TOKEN")

	// Create bot and enable debugging info
	// Note: Please keep in mind that default logger may expose sensitive information,
	// use in development only
	// (more on configuration in examples/configuration/main.go)
	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Get updates channel
	// (more on configuration in examples/updates_long_polling/main.go)
	updates, _ := bot.UpdatesViaLongPolling(context.Background(), nil)

	// Create bot handler and specify from where to get updates
	bh, _ := th.NewBotHandler(bot, updates)

	// Stop handling updates
	defer func() { _ = bh.Stop() }()

	// Register new handler with match on command `/start`
	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		// Send message

		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(

			tu.ID(update.Message.Chat.ID),
			fmt.Sprintf("Hello %s!", update.Message.From.FirstName),
		))
		return nil
	}, th.CommandEqual("hi"))

	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		// Send message
		fmt.Printf("New member %s", update.Message.NewChatMembers[0].Username)

		return nil
	}, NewMember())

	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		// Send message
		fmt.Printf("Left member %s", update.Message.LeftChatMember.Username)

		return nil
	}, LeftMember())

	// Loop through all updates when they came
	// for update := range updates {
	// 	if update.Message != nil {
	// 		//fmt.Printf("New message: %s", update.Message.Text)
	// 	}

	// 	if update.Message.LeftChatMember != nil {
	// 		fmt.Printf("Left member %s", update.Message.LeftChatMember.Username)
	// 	}

	// }
	//
	// // Start handling updates
	_ = bh.Start()
}
