package model

import "github.com/lib/pq"

type NewUser struct {
	UserId   int `gorm:"primaryKey";autoIncrement:false"`
	Username string
	Messages pq.Int64Array `gorm:"type:int[]"`
}
