package main

import (
	"time"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func reply(ctx *th.Context, chatId int64, replyTo int, text string) {
	_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
		tu.ID(chatId),
		text,
	).WithReplyParameters(&telego.ReplyParameters{MessageID: replyTo}))

}

func isDateValid(dateString string) bool {
	_, err := time.Parse("02/01/2006", dateString)
	return err == nil
}
