package models

import "time"

type Message struct {
	ID             uint         `gorm:"primarykey" json:"id"`
	ConversationID uint         `gorm:"not null" json:"conversation_id"`
	Conversation   Conversation `gorm:"foreignKey:ConversationID" json:"conversation"`
	UserID         uint         `gorm:"not null" json:"user_id"`
	User           User         `gorm:"foreignKey:UserID" json:"user"`
	Content        string       `gorm:"not null" json:"content"`
	MessageDate    time.Time    `gorm:"not null" json:"message_date"`
	PhotoURL       string
}
