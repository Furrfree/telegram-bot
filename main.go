package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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

	//bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	bot, err := telego.NewBot(botToken)

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

		if len(args) == 0 {
			reply(ctx, message.Chat.ID, message.MessageID, "Error: No se ha especificado el cumpleaños")
			return nil
		}

		birthdayDate := args[0]
		if !isDateValid(birthdayDate) {
			reply(ctx, message.Chat.ID, message.MessageID, "Error: El cumpleaños debe tener formato dd/mm/yyyy")
			return nil
		}
		date, _ := time.Parse("02/01/2006", birthdayDate)
		insertBirthday(
			db,
			message.From.ID,
			message.Chat.ID,
			date,
			message.From.Username,
		)

		reply(ctx, message.Chat.ID, message.MessageID, fmt.Sprintf("Añadido cumple de @%s el dia %s", message.From.Username, birthdayDate))
		return nil
	}, th.CommandEqual("add_cumple"))

	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		nextBirthday, _ := getNearestBirthday(db, message.Chat.ID)
		if nextBirthday == nil {
			reply(ctx, message.Chat.ID, message.MessageID, "No hay cumpleaños añadidos")
			return nil
		}

		reply(ctx, message.Chat.ID, message.MessageID, fmt.Sprintf("El siguiente cumple es el de @%s el dia %s", message.From.Username, nextBirthday.Date.Format("02/01/2006")))
		return nil
	}, th.CommandEqual("next_cumple"))

	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		// Send message
		fmt.Printf("New member %s", update.Message.NewChatMembers[0].Username)

		msg := sendMessage(ctx, update.Message.Chat.ID, fmt.Sprintf("Hola @%s", update.Message.NewChatMembers[0].Username))

		insertNewUser(db, update.Message.NewChatMembers[0].ID, msg.MessageID)

		return nil
	}, NewMember())

	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		// Send message
		fmt.Printf("Left member %s", update.Message.LeftChatMember.Username)

		welcomeMessageId := getWelcomeMessageId(db, update.Message.LeftChatMember.ID)

		err := ctx.Bot().DeleteMessage(ctx, &telego.DeleteMessageParams{
			ChatID:    update.Message.Chat.ChatID(),
			MessageID: welcomeMessageId,
		})

		if err != nil {
			fmt.Println(err)
		}

		deleteNewUser(db, welcomeMessageId)

		return nil
	}, LeftMember())

	privateChatCommands := telego.SetMyCommandsParams{
		Commands: []telego.BotCommand{
			{Command: "hi", Description: "Hello"},
		},
		Scope: tu.ScopeAllPrivateChats(),
	}

	groupCommands := telego.SetMyCommandsParams{
		Commands: []telego.BotCommand{
			{Command: "add_cumple", Description: "Añade tu cumpleaños al bot."},
			{Command: "next_cumple", Description: "Muestra el próximo cumpleaños"},
		},
		Scope: tu.ScopeAllGroupChats(),
	}

	bot.SetMyCommands(context.Background(), &privateChatCommands)
	bot.SetMyCommands(context.Background(), &groupCommands)

	_ = bh.Start()
}
