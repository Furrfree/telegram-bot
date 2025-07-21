package main

import (
	"context"
	"fmt"

	"github.com/furrfree/telegram-bot/commands"
	"github.com/furrfree/telegram-bot/configuration"
	"github.com/furrfree/telegram-bot/database"
	"github.com/furrfree/telegram-bot/logger"
	"github.com/furrfree/telegram-bot/service"
	"github.com/furrfree/telegram-bot/utils"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
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

	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		newMember := update.Message.NewChatMembers[0]
		fmt.Printf("Admission: New member %s", update.Message.NewChatMembers[0].Username)
		msg := utils.SendMarkdown(ctx, update.Message.Chat.ID, fmt.Sprintf(`
			¡Bienvenido/a, %s PARA ENTRAR:
			- Leer las [normas](%s) (y estar de acuerdo con ellas)
			- Ser mayor de edad: por las nuevas políticas de Telegram no podemos aceptar a personas menores de 18 años.
			- Presentarse: edad (obligatorio) de donde vienes, pronombres, nombres etc. Puedes usar esta [plantilla](%s)
			- Breve descripción y con qué podrías aportar (arte, quedadas, etc) (opcional)
			- Una vez os leamos seréis admitidos y entraréis en el grupo. Cuando entréis abandonad el grupo de admisión, por favor. Un saludo! 💜🐺
			`,
			update.Message.NewChatMembers[0].Username,
			configuration.ConfigInstance.RulesMessageUrl,
			configuration.ConfigInstance.PresentationTemplateMessageUrl))
		service.InsertNewUser(newMember.ID, newMember.Username, msg.MessageID)
		service.InsertNewUserMessage(newMember.ID, int64(msg.MessageID))
		return nil
	}, utils.NewMember(configuration.ConfigInstance.AdmissionGroupId))

	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		newUser := update.Message.NewChatMembers[0]
		fmt.Printf("Group: New member %s", newUser.Username)

		banError := bot.BanChatMember(ctx, &telego.BanChatMemberParams{
			ChatID: tu.ID(int64(configuration.ConfigInstance.AdmissionGroupId)),
			UserID: newUser.ID,
		})

		if banError != nil {
			logger.Error(banError)
			return nil
		}

		unbanError := bot.UnbanChatMember(ctx, &telego.UnbanChatMemberParams{
			ChatID: tu.ID(int64(configuration.ConfigInstance.AdmissionGroupId)),
			UserID: newUser.ID,
		})

		if unbanError != nil {
			logger.Error(unbanError)
			return nil
		}

		return nil
	}, utils.NewMember(configuration.ConfigInstance.GroupId))

	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		// Send message
		fmt.Printf("Left member %s", update.Message.LeftChatMember.Username)

		newUser := service.GetNewUserFromUserId(update.Message.LeftChatMember.ID)
		var messageIds []int

		for _, x := range service.GetNewUserFromUserId(int64(newUser.UserId)).Messages {
			messageIds = append(messageIds, int(x))
		}
		err := bot.DeleteMessages(ctx, &telego.DeleteMessagesParams{
			ChatID:     tu.ID(int64(configuration.ConfigInstance.AdmissionGroupId)),
			MessageIDs: messageIds,
		})

		if err != nil {
			logger.Error(err)
		}

		service.DeleteNewUser(newUser.UserId)

		return nil
	}, utils.LeftMember(configuration.ConfigInstance.AdmissionGroupId))

	_ = bh.Start()
}
