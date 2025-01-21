package models

type Conversation struct {
	ID       uint      `gorm:"primarykey" json:"id"`
	Messages []Message `gorm:"constraint:OnDelete:CASCADE" json:"messages"`
	Users    []*User   `gorm:"many2many:conversation_users;" json:"users"`
}
