package service

import (
	"context"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
)

func NewMember() telegohandler.Predicate {
	return func(ctx context.Context, update telego.Update) bool {
		return len(update.Message.NewChatMembers) != 0
	}
}

func LeftMember() telegohandler.Predicate {
	return func(ctx context.Context, update telego.Update) bool {
		return update.Message.LeftChatMember != nil
	}
}
