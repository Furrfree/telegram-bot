package commands

import (
	"fmt"

	"github.com/furrfree/telegram-bot/configuration"
	"github.com/furrfree/telegram-bot/service"
	"github.com/furrfree/telegram-bot/utils"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func admitir(bh *th.BotHandler, bot *telego.Bot) {
	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {

		if message.ReplyToMessage.From.ID == bot.ID() {
			utils.SendMessage(ctx, message.Chat.ChatID().ID, "Error: Responde con este comando al mensaje de presentación del usuario al que admitir")
		}
		newUser := service.GetNewUserFromUserId(message.ReplyToMessage.From.ID)
		service.InsertNewUserMessage(int64(newUser.UserId), int64(message.MessageID))

		if newUser.UserId == 0 {
			utils.SendMessage(ctx, int64(message.Chat.ID), "Error: No hay usuario que admitir")
			return nil
		}

		inviteLink, err := bot.CreateChatInviteLink(ctx, &telego.CreateChatInviteLinkParams{
			ChatID:      tu.ID(int64(configuration.Conf.GroupId)),
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
