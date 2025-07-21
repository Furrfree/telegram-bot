package commands

import (
	"fmt"
	"strings"

	"github.com/furrfree/telegram-bot/configuration"
	"github.com/furrfree/telegram-bot/logger"
	"github.com/furrfree/telegram-bot/service"
	"github.com/furrfree/telegram-bot/utils"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func admitir(bh *th.BotHandler, bot *telego.Bot) {
	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		_, _, args := tu.ParseCommand(message.Text)

		if len(args) == 0 {
			utils.Reply(ctx, message.Chat.ID, message.MessageID, "Error: No se ha especificado el usuario")
			return nil
		}

		username := strings.Split(args[0], "@")[1]
		logger.Log(username)

		newUser := service.GetNewUserByUsername(username)

		service.InsertNewUserMessage(int64(newUser.UserId), int64(message.MessageID))

		if newUser.UserId == 0 {
			utils.SendMessage(ctx, int64(message.Chat.ID), "Error: No hay usuario que admitir")
			return nil
		}

		inviteLink, err := bot.CreateChatInviteLink(ctx, &telego.CreateChatInviteLinkParams{
			ChatID:      tu.ID(int64(configuration.ConfigInstance.GroupId)),
			MemberLimit: 1,
		})

		if err != nil {
			fmt.Println(err)
			return nil
		}

		msg := utils.SendMessage(ctx, int64(message.Chat.ID), fmt.Sprintf("Aquí tienes el enlace al grupo %s. Una vez te unas se te echará de este grupo", inviteLink.InviteLink))
		service.InsertNewUserMessage(int64(newUser.UserId), int64(msg.MessageID))

		return nil
	}, th.CommandEqual("admitir"))

}
