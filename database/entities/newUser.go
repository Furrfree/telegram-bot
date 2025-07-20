package entities

type NewUser struct {
	UserId           int `gorm:"primaryKey";autoIncrement:false"`
	WelcomeMessageId int
}
