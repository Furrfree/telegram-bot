package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
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

	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	//bot, err := telego.NewBot(botToken)

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

		re := regexp.MustCompile(`^(0?[1-9]|[0-9]|3)/(0?[1-9]|1[0-2])/((19|20)\d{2})$`)

		if len(args) == 0 {
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
				tu.ID(message.Chat.ID),
				"Error: No se ha especificado el cumpleaños",
			).WithReplyParameters(&telego.ReplyParameters{MessageID: message.MessageID}))
			return nil
		}

		if !re.MatchString(args[0]) {
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
				tu.ID(message.Chat.ID),
				"Error: El cumpleaños debe tener formato dd/mm/yyyy",
			).WithReplyParameters(&telego.ReplyParameters{MessageID: message.MessageID}))
			return nil
		}

		userId := message.From.ID
		birthdayDate := args[0]
		groupId := message.Chat.ID
		format := "02/01/2006"
		date, _ := time.Parse(format, birthdayDate)

		insertBirthday(db, &entities.Birthday{
			UserId:   int(userId),
			Username: message.From.Username,
			GroupId:  int(groupId),
			Date:     date,
		})

		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(message.Chat.ID),
			fmt.Sprintf("Añadido cumple de @%s el dia %s", message.From.Username, birthdayDate),
		).WithReplyParameters(&telego.ReplyParameters{MessageID: message.MessageID}))

		return nil
	}, th.CommandEqual("add_cumple"))

	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		// Send message
		var nextBirthday entities.Birthday
		today := time.Now().Format("01-02") // Format as MM-DD
		db.Raw("SELECT * FROM birthdays WHERE group_id = ? ORDER BY strftime('%m-%d',date) >= strftime('%m-%d',datetime('now') ) DESC, strftime('%m-%d',date ) ASC LIMIT 1", message.Chat.ID, today).Scan(&nextBirthday)
		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(message.Chat.ID),
			fmt.Sprintf("El siguiente cumple es el de @%s el dia %s", message.From.Username, nextBirthday.Date.Format("02/01/2006")),
		).WithReplyParameters(&telego.ReplyParameters{MessageID: message.MessageID}))

		return nil
	}, th.CommandEqual("next_cumple"))

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

	groupCommands := telego.SetMyCommandsParams{
		Commands: []telego.BotCommand{
			{Command: "add_cumple", Description: "Añade tu cumpleaños al bot."},
			{Command: "next_cumple", Description: "Muestra el próximo cumpleaños"},
		},
		Scope:        tu.ScopeAllGroupChats(),
		LanguageCode: "es",
	}

	bot.SetMyCommands(context.Background(), &privateChatCommands)
	bot.SetMyCommands(context.Background(), &groupCommands)

	_ = bh.Start()
}
