package models

import "time"

type Message struct {
	ID             uint      `gorm:"primarykey"`
	ConversationID uint      `gorm:"not null"`
	UserID         uint      `gorm:"not null"`
	Content        string    `gorm:"not null"`
	MessageDate    time.Time `gorm:"not null"`
	PhotoURL       string
}
