package main

import (
	"context"
	"fmt"
	"time"

	"github.com/furrfree/telegram-bot/logger"
	"github.com/furrfree/telegram-bot/service"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func main() {
	// Set up DB
	service.InitializeDb()
	config := service.GetConfig()

	//bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	bot, err := telego.NewBot(config.Token)

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

	// Register new handler with match on command `/start`
	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		// Send message

		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(

			tu.ID(update.Message.Chat.ID),
			fmt.Sprintf("Hello %s!", update.Message.From.FirstName),
		))
		return nil
	}, th.CommandEqual("hi"))

	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		// Send message
		_, _, args := tu.ParseCommand(message.Text)

		if len(args) == 0 {
			service.Reply(ctx, message.Chat.ID, message.MessageID, "Error: No se ha especificado el cumpleaños")
			return nil
		}

		birthdayDate := args[0]
		if !service.IsDateValid(birthdayDate) {
			service.Reply(ctx, message.Chat.ID, message.MessageID, "Error: El cumpleaños debe tener formato dd/mm/yyyy")
			return nil
		}
		date, _ := time.Parse("02/01/2006", birthdayDate)
		service.InsertBirthday(
			message.From.ID,
			message.Chat.ID,
			date,
			message.From.Username,
		)

		service.Reply(ctx, message.Chat.ID, message.MessageID, fmt.Sprintf("Añadido cumple de @%s el dia %s", message.From.Username, birthdayDate))
		return nil
	}, th.CommandEqual("add_cumple"))

	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		nextBirthday, _ := service.GetNearestBirthday(message.Chat.ID)
		if nextBirthday == nil {
			service.Reply(ctx, message.Chat.ID, message.MessageID, "No hay cumpleaños añadidos")
			return nil
		}

		service.Reply(ctx, message.Chat.ID, message.MessageID, fmt.Sprintf("El siguiente cumple es el de @%s el dia %s", message.From.Username, nextBirthday.Date.Format("02/01/2006")))
		return nil
	}, th.CommandEqual("next_cumple"))

	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {

		// TODO: Change to get the replied message id
		welcomeMessageId := message.ReplyToMessage.MessageID
		newUser := service.GetNewUserByMessageId(int64(welcomeMessageId))

		if newUser.UserId == 0 {
			service.SendMessage(ctx, int64(message.Chat.ID), "Error: No hay usuario que admitir")
			return nil
		}

		inviteLink, err := bot.CreateChatInviteLink(ctx, &telego.CreateChatInviteLinkParams{
			ChatID:      tu.ID(int64(config.GroupId)),
			MemberLimit: 1,
		})

		if err != nil {
			fmt.Println(err)
			return nil
		}
		service.SendMessage(ctx, int64(message.Chat.ID), fmt.Sprintf("Aquí tienes el enlace al grupo %s. Una vez te unas se te echará de este grupo", inviteLink.InviteLink))

		return nil
	}, th.CommandEqual("admitir"))

	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		fmt.Printf("Admission: New member %s", update.Message.NewChatMembers[0].Username)
		msg := service.SendMarkdown(ctx, update.Message.Chat.ID, fmt.Sprintf(`
			¡Bienvenido/a, %s PARA ENTRAR:
			- Leer las [normas](%s) (y estar de acuerdo con ellas)
			- Ser mayor de edad: por las nuevas políticas de Telegram no podemos aceptar a personas menores de 18 años.
			- Presentarse: edad (obligatorio) de donde vienes, pronombres, nombres etc. Puedes usar esta [plantilla](%s)
			- Breve descripción y con qué podrías aportar (arte, quedadas, etc) (opcional)
			- Una vez os leamos seréis admitidos y entraréis en el grupo. Cuando entréis abandonad el grupo de admisión, por favor. Un saludo! 💜🐺
			`,
			update.Message.NewChatMembers[0].Username,
			config.RulesMessageUrl,
			config.PresentationTemplateMessageUrl))
		service.InsertNewUser(update.Message.NewChatMembers[0].ID, msg.MessageID)
		return nil
	}, service.NewMember(config.AdmissionGroupId))

	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		newUser := update.Message.NewChatMembers[0]
		fmt.Printf("Group: New member %s", newUser.Username)

		banError := bot.BanChatMember(ctx, &telego.BanChatMemberParams{
			ChatID: tu.ID(int64(config.AdmissionGroupId)),
			UserID: newUser.ID,
		})

		if banError != nil {
			logger.Error(banError)
			return nil
		}

		unbanError := bot.UnbanChatMember(ctx, &telego.UnbanChatMemberParams{
			ChatID: tu.ID(int64(config.AdmissionGroupId)),
			UserID: newUser.ID,
		})

		if unbanError != nil {
			logger.Error(unbanError)
			return nil
		}

		return nil
	}, service.NewMember(config.GroupId))

	bh.Handle(func(ctx *th.Context, update telego.Update) error {
		// Send message
		fmt.Printf("Left member %s", update.Message.LeftChatMember.Username)

		newUser := service.GetNewUserFromUserId(update.Message.LeftChatMember.ID)

		err := ctx.Bot().DeleteMessage(ctx, &telego.DeleteMessageParams{
			ChatID:    update.Message.Chat.ChatID(),
			MessageID: newUser.WelcomeMessageId,
		})

		if err != nil {
			logger.Error(err)
		}

		service.DeleteNewUser(newUser.UserId)

		return nil
	}, service.LeftMember())

	privateChatCommands := telego.SetMyCommandsParams{
		Commands: []telego.BotCommand{
			{Command: "hi", Description: "Hello"},
		},
		Scope: tu.ScopeAllPrivateChats(),
	}

	groupCommands := telego.SetMyCommandsParams{
		Commands: []telego.BotCommand{
			{Command: "add_cumple", Description: "Añade tu cumpleaños al bot."},
			{Command: "next_cumple", Description: "Muestra el próximo cumpleaños"},
		},
		Scope: tu.ScopeAllGroupChats(),
	}

	admissionGroupCommands := telego.SetMyCommandsParams{
		Commands: []telego.BotCommand{
			{Command: "admitir", Description: "Admite a un usuario"},
		},
		Scope: tu.ScopeChatAdministrators(telego.ChatID{ID: int64(config.AdmissionGroupId)}),
	}

	bot.SetMyCommands(context.Background(), &privateChatCommands)
	bot.SetMyCommands(context.Background(), &groupCommands)
	bot.SetMyCommands(context.Background(), &admissionGroupCommands)

	_ = bh.Start()
}
