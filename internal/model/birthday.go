package model

import (
	"time"
)

type Birthday struct {
	UserId   int       `gorm:"primaryKey";autoIncrement:false"`
	GroupId  int       `gorm:"primaryKey";autoIncrement:false"`
	Date     time.Time `gorm:"type:DATE;"`
	Username string
}
