// internal/shared/repositories/conversation_repository.go
package repositories

import (
	"socialmedia-backend/internal/shared/db"
	"socialmedia-backend/internal/shared/models"

	"gorm.io/gorm"
)

type ConversationRepository struct {
	db *gorm.DB
}

func NewConversationRepository() *ConversationRepository {
	return &ConversationRepository{
		db: db.GetDB(),
	}
}

func (r *ConversationRepository) Create(conv *models.Conversation) error {
	return r.db.Create(conv).Error
}

func (r *ConversationRepository) GetByID(id uint) (*models.Conversation, error) {
	var conv models.Conversation
	err := r.db.Preload("Messages").Preload("Users").First(&conv, id).Error
	return &conv, err
}

func (r *ConversationRepository) GetUserConversations(userID uint) ([]models.Conversation, error) {
	var conversations []models.Conversation
	err := r.db.Joins("JOIN conversation_users ON conversations.id = conversation_users.conversation_id").
		Where("conversation_users.user_id = ?", userID).
		Preload("Messages").
		Preload("Users").
		Find(&conversations).Error
	return conversations, err
}
