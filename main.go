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
		userId := message.From.ID
		birthdayDate := args[0]
		groupId := message.Chat.ID
		format := "16-01-2001"
		date, _ := time.Parse(format, birthdayDate)

		db.Create(&entities.Birthday{
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
		//groupId := message.Chat.ID
		var nextCumple entities.Birthday
		//var result entities.Birthday

		var birthdays []entities.Birthday
		println("test")

		//date1, _ := time.Parse("11/11/2000", "04/09/1997")
		date2, _ := time.Parse("02/01/2006", "23/10/2015")

		db.Create(&entities.Birthday{
			UserId:   1,
			GroupId:  int(message.Chat.ID),
			Date:     time.Now(),
			Username: "test",
		})
		db.Create(&entities.Birthday{
			UserId:   2,
			GroupId:  int(message.Chat.ID),
			Date:     date2,
			Username: "testoo",
		})

		db.Find(&birthdays)

		for _, x := range birthdays {
			fmt.Println(x.GroupId, x.UserId, x.Username, x.Date)
		}

		db.First(&nextCumple)
		var test entities.Birthday
		today := time.Now().Format("01-02") // Format as MM-DD
		db.Raw(`
    SELECT *
    FROM birthdays
    WHERE group_id = ?
    ORDER BY
        strftime('%m-%d', date) >= ? DESC,
        ABS(julianday(date) - julianday('now'))
    LIMIT 1
`, 123, today).Scan(&test)

		fmt.Println(test.Date)

		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(message.Chat.ID),
			fmt.Sprintf("Añadido cumple de @%s el dia %s", message.From.Username, "safd"),
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

	bot.SetMyCommands(context.Background(), &privateChatCommands)

	_ = bh.Start()
}
