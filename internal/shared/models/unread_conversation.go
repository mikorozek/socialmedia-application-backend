package models

import "time"

type UnreadConversation struct {
	ID                 uint         `gorm:"primarykey" json:"id"`
	ConversationID     uint         `gorm:"not null" json:"conversation_id"`
	Conversation       Conversation `gorm:"foreignKey:ConversationID" json:"conversation"`
	UserID             uint         `gorm:"not null" json:"user_id"`
	User               User         `gorm:"foreignKey:UserID" json:"user"`
	LastMessageContent string       `gorm:"not null" json:"last_message_content"`
	LastMessageTime    time.Time    `gorm:"not null" json:"last_message_time"`
}
