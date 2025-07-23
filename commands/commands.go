package commands

import (
	"context"
	"fmt"

	"github.com/furrfree/telegram-bot/configuration"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func AddCommands(bh *th.BotHandler, bot *telego.Bot) {

	fmt.Println(configuration.Conf.AdmissionGroupId)
	// Private chat commands
	addPrivateCommands(bh, bot)
	addGroupCommands(bh, bot)
	addGroupAdminCommands(bh, bot)
}

func addPrivateCommands(bh *th.BotHandler, bot *telego.Bot) {
	hi(bh)

	var PrivateChatCommands = telego.SetMyCommandsParams{
		Commands: []telego.BotCommand{
			{Command: "hi", Description: "Hello"},
		},
		Scope: tu.ScopeAllPrivateChats(),
	}
	bot.SetMyCommands(context.Background(), &PrivateChatCommands)
}

func addGroupCommands(bh *th.BotHandler, bot *telego.Bot) {
	add_cumple(bh)
	next_cumple(bh)

	var GroupCommands = telego.SetMyCommandsParams{
		Commands: []telego.BotCommand{
			{Command: "add_cumple", Description: "Añade tu cumpleaños al bot."},
			{Command: "next_cumple", Description: "Muestra el próximo cumpleaños"},
		},
		Scope: tu.ScopeAllGroupChats(),
	}
	bot.SetMyCommands(context.Background(), &GroupCommands)
}

func addGroupAdminCommands(bh *th.BotHandler, bot *telego.Bot) {
	admitir(bh, bot)
	var AdmissionGroupCommands = telego.SetMyCommandsParams{
		Commands: []telego.BotCommand{
			{Command: "admitir", Description: "Admite a un usuario"},
		},
		Scope: tu.ScopeChatAdministrators(telego.ChatID{ID: int64(configuration.Conf.AdmissionGroupId)}),
	}
	bot.SetMyCommands(context.Background(), &AdmissionGroupCommands)
}
