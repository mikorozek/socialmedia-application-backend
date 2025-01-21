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

func (r *MessageRepository) AddMessage(message *models.Message) error {
	return r.db.Create(message).Error
}

func (r *MessageRepository) GetMessages(conversationID uint, limit int, offset int) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.Where("conversation_id = ?", conversationID).
		Order("message_date DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	return messages, err
}

func (r *MessageRepository) GetMessageByID(messageID uint) (*models.Message, error) {
	var message models.Message
	err := r.db.First(&message, messageID).Error
	return &message, err
}

func (r *MessageRepository) UpdateMessage(message *models.Message) error {
	return r.db.Save(message).Error
}

func (r *MessageRepository) DeleteMessage(messageID uint) error {
	return r.db.Delete(&models.Message{}, messageID).Error
}
