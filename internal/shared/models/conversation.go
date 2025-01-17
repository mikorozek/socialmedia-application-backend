package models

type Conversation struct {
	ID       uint `gorm:"primarykey"`
	Messages []Message
	Users    []User `gorm:"many2many:conversation_users;"`
}
