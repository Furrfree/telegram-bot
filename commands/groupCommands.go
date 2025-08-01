package commands

import (
	"fmt"
	"time"

	"github.com/furrfree/telegram-bot/service"
	"github.com/furrfree/telegram-bot/utils"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func add_cumple(bh *th.BotHandler) {
	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		_, _, args := tu.ParseCommand(message.Text)

		if len(args) == 0 {
			utils.Reply(ctx, message.Chat.ID, message.MessageID, "Error: No se ha especificado el cumpleaños")
			return nil
		}

		birthdayDate := args[0]
		if !utils.IsDateValid(birthdayDate) {
			utils.Reply(ctx, message.Chat.ID, message.MessageID, "Error: El cumpleaños debe tener formato dd/mm/yyyy")
			return nil
		}
		date, _ := time.Parse("02/01/2006", birthdayDate)
		service.InsertBirthday(
			message.From.ID,
			message.Chat.ID,
			date,
			message.From.Username,
		)

		utils.Reply(ctx, message.Chat.ID, message.MessageID, fmt.Sprintf("Añadido cumple de @%s el dia %s", message.From.Username, birthdayDate))
		return nil
	}, th.CommandEqual("add_cumple"))

}

func next_cumple(bh *th.BotHandler) {
	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		nextBirthday, _ := service.GetNearestBirthday(message.Chat.ID)
		if nextBirthday == nil {
			utils.Reply(ctx, message.Chat.ID, message.MessageID, "No hay cumpleaños añadidos")
			return nil
		}

		utils.Reply(ctx, message.Chat.ID, message.MessageID, fmt.Sprintf("El siguiente cumple es el de @%s el %s", nextBirthday.Username, utils.ParseBirthday(nextBirthday.Date)))
		return nil
	}, th.CommandEqual("next_cumple"))

}
