package models

type User struct {
	ID           uint            `gorm:"primarykey" json:"id"`
	Username     string          `gorm:"unique;not null" json:"username"`
	Email        string          `gorm:"unique;not null" json:"email"`
	PasswordHash string          `gorm:"not null" json:"password_hash"`
	Conversation []*Conversation `gorm:"many2many:conversation_user;constraint:OnDelete:CASCADE" json:"conversation"`
	Messages     []Message       `gorm:"constraint:OnDelete:CASCADE" json:"messages"`
}
