package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/furrfree/telegram-bot/entities"
	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func main() {

	// Set up DB
	db := setupDb()
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

	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		// Send message
		_, _, args := tu.ParseCommand(message.Text)
		userId := message.From.ID
		birthdayDate := args[0]
		format := "16-01-2001"
		date, _ := time.Parse(format, birthdayDate)

		db.Create(&entities.Birthday{
			UserId:   int(userId),
			Username: message.From.Username,
			Date:     &date,
		})

		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(message.Chat.ID),
			fmt.Sprintf("Añadido cumple de @%s el dia %s", message.From.Username, birthdayDate),
		).WithReplyParameters(&telego.ReplyParameters{MessageID: message.MessageID}))

		return nil
	}, th.CommandEqual("add_cumple"))

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

	privateChatCommands := telego.SetMyCommandsParams{
		Commands: []telego.BotCommand{
			{Command: "Hi", Description: "Hello"},
		},
		Scope:        tu.ScopeAllPrivateChats(),
		LanguageCode: "es",
	}

	bot.SetMyCommands(context.Background(), &privateChatCommands)

	_ = bh.Start()
}
