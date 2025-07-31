package model

import (
	"time"

	"github.com/lib/pq"
)

type NewUser struct {
	UserId     int `gorm:"primaryKey";autoIncrement:false"`
	ChatID     int
	Username   string
	Messages   pq.Int64Array `gorm:"type:int[]"`
	DateJoined time.Time
}
