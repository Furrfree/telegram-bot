package main

import (
	"context"

	"github.com/furrfree/telegram-bot/commands"
	"github.com/furrfree/telegram-bot/configuration"
	"github.com/furrfree/telegram-bot/coroutines"
	"github.com/furrfree/telegram-bot/database"
	"github.com/furrfree/telegram-bot/handlers"
	"github.com/furrfree/telegram-bot/logger"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

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
	coroutines.AddCoroutines(bot)
	logger.Log("Bot started")
	bh.Start()

}
