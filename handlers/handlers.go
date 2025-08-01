package handlers

import (
	"fmt"

	"github.com/furrfree/telegram-bot/configuration"
	"github.com/furrfree/telegram-bot/coroutines"
	"github.com/furrfree/telegram-bot/logger"
	"github.com/furrfree/telegram-bot/service"
	"github.com/furrfree/telegram-bot/utils"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func AddHandlers(bh *th.BotHandler, bot *telego.Bot) {
	newMemberAdmissionGroup(bh, bot)
	newGroupMember(bh, bot)
	leaveAdmissionGroup(bh, bot)

}

func newMemberAdmissionGroup(bh *th.BotHandler, bot *telego.Bot) {
	bh.Handle(func(ctx *th.Context, update telego.Update) error {

		// Remove join message
		err := bot.DeleteMessage(ctx, &telego.DeleteMessageParams{
			ChatID:    tu.ID(int64(configuration.Conf.AdmissionGroupId)),
			MessageID: update.Message.MessageID,
		})

		if err != nil {
			logger.Error("Could not delete user joined message")
		}

		newMember := update.Message.NewChatMembers[0]
		logger.Log(fmt.Sprintf("Admission: New member %s", update.Message.NewChatMembers[0].Username))
		msg := utils.SendMarkdown(ctx, update.Message.Chat.ID, fmt.Sprintf(`
			¡Bienvenido/a, %s PARA ENTRAR:
			- Leer las [normas](%s) (y estar de acuerdo con ellas)
			- Ser mayor de edad: por las nuevas políticas de Telegram no podemos aceptar a personas menores de 18 años.
			- Presentarse: edad (obligatorio) de donde vienes, pronombres, nombres etc. Puedes usar esta [plantilla](%s)
			- Breve descripción y con qué podrías aportar (arte, quedadas, etc) (opcional)
			- Una vez os leamos seréis admitidos y entraréis en el grupo. Cuando entréis abandonad el grupo de admisión, por favor. Un saludo! 💜🐺
			`,
			update.Message.NewChatMembers[0].Username,
			configuration.Conf.RulesMessageUrl,
			configuration.Conf.PresentationTemplateMessageUrl))
		service.InsertNewUser(newMember.ID, update.Message.Chat.ID, newMember.Username)
		service.InsertNewUserMessage(newMember.ID, int64(msg.MessageID))
		// Send signal to couroutine of new member in admission
		select {
		case coroutines.NewUserAddedChan <- struct{}{}:
		default:
		}

		return nil
	}, utils.NewMember(configuration.Conf.AdmissionGroupId))
}

func newGroupMember(bh *th.BotHandler, bot *telego.Bot) {
	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		newUser := update.Message.NewChatMembers[0]
		logger.Log(fmt.Sprintf("Group: New member %s", newUser.Username))

		if banError := bot.BanChatMember(ctx, &telego.BanChatMemberParams{
			ChatID: tu.ID(int64(configuration.Conf.AdmissionGroupId)),
			UserID: newUser.ID,
		}); banError != nil {
			logger.Error(banError)
			return nil
		}

		if unbanError := bot.UnbanChatMember(ctx, &telego.UnbanChatMemberParams{
			ChatID: tu.ID(int64(configuration.Conf.AdmissionGroupId)),
			UserID: newUser.ID,
		}); unbanError != nil {
			logger.Error(unbanError)
			return nil
		}

		return nil
	}, utils.NewMember(configuration.Conf.GroupId))
}

func leaveAdmissionGroup(bh *th.BotHandler, bot *telego.Bot) {
	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		logger.Log(fmt.Sprintf("Admission: Left member %s", update.Message.LeftChatMember.Username))

		// Remove left message
		if errDeleteingLeftMessage := bot.DeleteMessage(ctx, &telego.DeleteMessageParams{
			ChatID:    tu.ID(int64(configuration.Conf.AdmissionGroupId)),
			MessageID: update.Message.MessageID,
		}); errDeleteingLeftMessage != nil {
			logger.Error("Could not delete user left message")
		}

		newUser := service.GetNewUserByUsername(update.Message.LeftChatMember.Username)

		// Get messages
		var messageIds []int
		for _, x := range newUser.Messages {
			messageIds = append(messageIds, int(x))
		}

		// Remove messages
		if err := bot.DeleteMessages(ctx, &telego.DeleteMessagesParams{
			ChatID:     tu.ID(int64(configuration.Conf.AdmissionGroupId)),
			MessageIDs: messageIds,
		}); err != nil {
			logger.Error(err)
		}
		service.DeleteNewUser(newUser.UserId)
		return nil
	}, utils.LeftMember(configuration.Conf.AdmissionGroupId))
}
