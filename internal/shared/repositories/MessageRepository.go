package repositories

import (
	"socialmedia-backend/internal/shared/db"
	"socialmedia-backend/internal/shared/models"

	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository() *MessageRepository {
	return &MessageRepository{
		db: db.GetDB(),
	}
}

func (r *ConversationRepository) AddMessage(message *models.Message) error {
	return r.db.Create(message).Error
}

func (r *ConversationRepository) GetMessages(conversationID uint, limit int, offset int) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.Where("conversation_id = ?", conversationID).
		Order("message_date DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	return messages, err
}

func (r *ConversationRepository) GetMessageByID(messageID uint) (*models.Message, error) {
	var message models.Message
	err := r.db.First(&message, messageID).Error
	return &message, err
}

func (r *ConversationRepository) UpdateMessage(message *models.Message) error {
	return r.db.Save(message).Error
}

func (r *ConversationRepository) DeleteMessage(messageID uint) error {
	return r.db.Delete(&models.Message{}, messageID).Error
}

func (r *ConversationRepository) AddUnreadMessage(messageID uint, userID uint, conversationID uint) error {
	unreadMessage := &models.UnreadMessage{
		MessageID:      messageID,
		UserID:         userID,
		ConversationID: conversationID,
	}
	return r.db.Create(unreadMessage).Error
}

func (r *ConversationRepository) DeleteUnreadMessage(messageID uint, userID uint, conversationID uint) error {
	return r.db.Where("message_id = ? AND user_id = ? AND conversation_id = ?",
		messageID, userID, conversationID).
		Delete(&models.UnreadMessage{}).Error
}

func (r *ConversationRepository) DeleteUnreadMessages(messageID uint) error {
	return r.db.Where("message_id = ?", messageID).Delete(&models.UnreadMessage{}).Error
}

func (r *ConversationRepository) GetUnreadMessagesForUser(userID uint) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.Table("messages").
		Joins("JOIN unread_messages ON messages.id = unread_messages.message_id").
		Where("unread_messages.user_id = ?", userID).
		Find(&messages).Error
	return messages, err
}
