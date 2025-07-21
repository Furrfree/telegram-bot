package service

import (
	"context"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/mymmrac/telego/telegoutil"
)

func NewMember(chatId int) telegohandler.Predicate {
	return func(ctx context.Context, update telego.Update) bool {
		return len(update.Message.NewChatMembers) != 0 && update.Message.Chat.ChatID() == telegoutil.ID(int64(chatId))
	}
}

func LeftMember() telegohandler.Predicate {
	return func(ctx context.Context, update telego.Update) bool {
		return update.Message.LeftChatMember != nil
	}
}
