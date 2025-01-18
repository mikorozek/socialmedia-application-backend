package models

type UnreadMessage struct {
	MessageID      uint         `gorm:"primarykey"`
	Message        Message      `gorm:"foreignKey:MessageID"`
	UserID         uint         `gorm:"primarykey"`
	User           User         `gorm:"foreignKey:UserID"`
	ConversationID uint         `gorm:"not null"`
	Conversation   Conversation `gorm:"foreignKey:ConversationID"`
}
