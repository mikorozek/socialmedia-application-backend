package models

import "time"

type Message struct {
	ID             uint         `gorm:"primarykey"`
	ConversationID uint         `gorm:"not null"`
	Conversation   Conversation `gorm:"foreignKey:ConversationID"`
	UserID         uint         `gorm:"not null"`
	User           User         `gorm:"foreignKey:UserID"`
	Content        string       `gorm:"not null"`
	MessageDate    time.Time    `gorm:"not null"`
	PhotoURL       string
}
