package main

import (
	"context"

	"github.com/furrfree/telegram-bot/internal/commands"
	"github.com/furrfree/telegram-bot/internal/configuration"
	"github.com/furrfree/telegram-bot/internal/coroutines"
	"github.com/furrfree/telegram-bot/internal/database"
	"github.com/furrfree/telegram-bot/internal/handlers"
	"github.com/furrfree/telegram-bot/internal/logger"
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
