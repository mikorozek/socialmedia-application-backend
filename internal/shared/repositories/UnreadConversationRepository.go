package repositories

import (
	"socialmedia-backend/internal/shared/db"
	"socialmedia-backend/internal/shared/models"
	"time"

	"gorm.io/gorm"
)

type UnreadConversationRepository struct {
	db *gorm.DB
}

func NewUnreadConversationRepository() *UnreadConversationRepository {
	return &UnreadConversationRepository{
		db: db.GetDB(),
	}
}

func (r *UnreadConversationRepository) UpdateUnreadConversation(conversationID uint, userID uint, content string, messageTime time.Time) error {
	return r.db.Where(models.UnreadConversation{
		ConversationID: conversationID,
		UserID:         userID,
	}).Assign(models.UnreadConversation{
		LastMessageContent: content,
		LastMessageTime:    messageTime,
	}).FirstOrCreate(&models.UnreadConversation{}).Error
}

func (r *UnreadConversationRepository) GetUnreadConversations(userID uint) ([]models.UnreadConversation, error) {
	var unreadConvs []models.UnreadConversation
	err := r.db.Where("user_id = ?", userID).
		Preload("Conversation").
		Preload("Conversation.Users").
		Order("last_message_time DESC").
		Find(&unreadConvs).Error
	return unreadConvs, err
}

func (r *UnreadConversationRepository) MarkAsRead(conversationID uint, userID uint) error {
	return r.db.Where("conversation_id = ? AND user_id = ?",
		conversationID, userID).
		Delete(&models.UnreadConversation{}).Error
}
