package entities

import (
	"time"

	"gorm.io/gorm"
)

type Birthday struct {
	gorm.Model
	UserId   int
	Date     *time.Time
	Username string
}
