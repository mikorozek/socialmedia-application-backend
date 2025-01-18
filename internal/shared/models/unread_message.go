package models

import "time"

type UnreadConversation struct {
	ID                 uint         `gorm:"primarykey"`
	ConversationID     uint         `gorm:"not null"`
	Conversation       Conversation `gorm:"foreignKey:ConversationID"`
	UserID             uint         `gorm:"not null"`
	User               User         `gorm:"foreignKey:UserID"`
	LastMessageContent string       `gorm:"not null"`
	LastMessageTime    time.Time    `gorm:"not null"`
}
