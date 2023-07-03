package database

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           uint `gorm:"primaryKey"`
	Username     string
	PasswordHash string
	Email        string
	Polls        []Poll `gorm:"foreignKey:CreatedByID"`
	Votes        []Vote `gorm:"foreignKey:UserID"`
}

type Poll struct {
	gorm.Model
	ID          uint `gorm:"primaryKey"`
	Title       string
	CreatedByID uint
	CreatedBy   User       `gorm:"foreignKey:CreatedByID"`
	Questions   []Question `gorm:"foreignKey:PollID"`
}

type Question struct {
	gorm.Model
	ID      uint `gorm:"primaryKey"`
	Text    string
	PollID  uint
	Poll    Poll     `gorm:"foreignKey:PollID"`
	Options []Option `gorm:"foreignKey:QuestionID"`
}

type Option struct {
	gorm.Model
	ID         uint `gorm:"primaryKey"`
	Text       string
	QuestionID uint
	Question   Question `gorm:"foreignKey:QuestionID"`
	Votes      []Vote   `gorm:"foreignKey:OptionID"`
}

type Vote struct {
	gorm.Model
	ID       uint `gorm:"primaryKey"` // This field becomes the primary key
	UserID   uint `gorm:"foreignKey:UserID"`
	User     User
	OptionID uint `gorm:"foreignKey:OptionID"`
	Option   Option
}
