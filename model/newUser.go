package model

type NewUser struct {
	UserId           int `gorm:"primaryKey";autoIncrement:false"`
	WelcomeMessageId int
}
