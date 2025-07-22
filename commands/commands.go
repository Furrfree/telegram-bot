package commands

import (
	"context"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func AddCommands(bh *th.BotHandler, bot *telego.Bot) {
	// Private chat commands
	addPrivateCommands(bh, bot)
	addGroupCommands(bh, bot)
	addGroupAdminCommands(bh, bot)
}

func addPrivateCommands(bh *th.BotHandler, bot *telego.Bot) {
	hi(bh)

	bot.SetMyCommands(context.Background(), &PrivateChatCommands)
}

func addGroupCommands(bh *th.BotHandler, bot *telego.Bot) {
	add_cumple(bh)
	next_cumple(bh)

	bot.SetMyCommands(context.Background(), &GroupCommands)
}

func addGroupAdminCommands(bh *th.BotHandler, bot *telego.Bot) {
	admitir(bh, bot)

	bot.SetMyCommands(context.Background(), &AdmissionGroupCommands)
}
