package repositories

import (
	"socialmedia-backend/internal/shared/db"
	"socialmedia-backend/internal/shared/models"

	"gorm.io/gorm"
)

type UnreadMessageRepository struct {
	db *gorm.DB
}

func NewUnreadMessageRepository() *UnreadMessageRepository {
	return &UnreadMessageRepository{
		db: db.GetDB(),
	}
}

func (r *UnreadMessageRepository) AddUnreadMessage(messageID uint, userID uint, conversationID uint) error {
	unreadMessage := &models.UnreadMessage{
		MessageID:      messageID,
		UserID:         userID,
		ConversationID: conversationID,
	}
	return r.db.Create(unreadMessage).Error
}

func (r *UnreadMessageRepository) DeleteUnreadMessage(messageID uint, userID uint, conversationID uint) error {
	return r.db.Where("message_id = ? AND user_id = ? AND conversation_id = ?",
		messageID, userID, conversationID).
		Delete(&models.UnreadMessage{}).Error
}
