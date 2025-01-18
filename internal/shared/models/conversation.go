package models

type Conversation struct {
	ID       uint      `gorm:"primarykey"`
	Messages []Message `gorm:"constraint:OnDelete:CASCADE"`
	Users    []User    `gorm:"many2many:conversation_users;"`
}
