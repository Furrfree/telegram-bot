package main

import (
	"context"

	"github.com/furrfree/telegram-bot/commands"
	"github.com/furrfree/telegram-bot/configuration"
	"github.com/furrfree/telegram-bot/database"
	"github.com/furrfree/telegram-bot/handlers"
	"github.com/furrfree/telegram-bot/logger"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func main() {
	// Set up DB
	database.InitializeDb()
	configuration.InitializeConfig()

	//bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	bot, err := telego.NewBot(configuration.ConfigInstance.Token)

	if err != nil {
		logger.Fatal(err)
	}

	// Get updates channel
	// (more on configuration in examples/updates_long_polling/main.go)
	updates, _ := bot.UpdatesViaLongPolling(context.Background(), nil)

	// Create bot handler and specify from where to get updates
	bh, _ := th.NewBotHandler(bot, updates)

	// Stop handling updates
	defer func() { _ = bh.Stop() }()

	commands.AddCommands(bh, bot)
	handlers.AddHandlers(bh, bot)

	_ = bh.Start()
}
