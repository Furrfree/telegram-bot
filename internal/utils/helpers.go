package utils

import (
	"time"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func Reply(ctx *th.Context, chatId int64, replyTo int, text string) {
	_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
		tu.ID(chatId),
		text,
	).WithReplyParameters(&telego.ReplyParameters{MessageID: replyTo}))
}

func SendMessage(ctx *th.Context, chatId int64, text string) *telego.Message {
	msg, _ := ctx.Bot().SendMessage(ctx, tu.Message(
		tu.ID(chatId),
		text,
	))
	return msg
}

func SendMarkdown(ctx *th.Context, chatId int64, text string) *telego.Message {
	msg, _ := ctx.Bot().SendMessage(ctx, tu.Message(
		tu.ID(chatId),
		text,
	).WithParseMode(telego.ModeMarkdown))
	return msg
}

func IsDateValid(dateString string) bool {
	_, err := time.Parse("02/01/2006", dateString)
	return err == nil
}
