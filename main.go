package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
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

	// Loop through all updates when they came
	for update := range updates {
		if update.Message != nil {
			//fmt.Printf("New message: %s", update.Message.Text)
		}

		if len(update.Message.NewChatMembers) != 0 {
			for i, newUser := range update.Message.NewChatMembers {
				fmt.Printf("New member %d %s", i, newUser.Username)
			}
		}

		if update.Message.LeftChatMember != nil {
			fmt.Printf("Left member %s", update.Message.LeftChatMember.Username)
		}

	}
}
