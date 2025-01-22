package repositories

import (
	"socialmedia-backend/internal/shared/db"
	"socialmedia-backend/internal/shared/models"
	"sort"

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

type ConversationWithLastMessage struct {
	Conversation models.Conversation
	LastMessage  *models.Message
	IsUnread     bool
}

func (r *ConversationRepository) Create(conv *models.Conversation) error {
	return r.db.Create(conv).Error
}

func (r *ConversationRepository) AddUserToConversation(convID uint, userID uint) error {
	return r.db.Exec("INSERT INTO conversation_users (conversation_id, user_id) VALUES (?, ?)", convID, userID).Error
}

func (r *ConversationRepository) GetByID(convID uint) (*models.Conversation, error) {
	var conv models.Conversation
	err := r.db.Preload("Users").First(&conv, convID).Error
	return &conv, err
}

func (r *ConversationRepository) GetUserConversations(userID uint) ([]models.Conversation, error) {
	var conversations []models.Conversation
	err := r.db.
		Joins("JOIN conversation_users ON conversations.id = conversation_users.conversation_id").
		Where("conversation_users.user_id = ?", userID).
		Preload("Users").
		Find(&conversations).Error
	return conversations, err
}

func (r *ConversationRepository) CheckUserInConversation(convID uint, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Conversation{}).
		Joins("JOIN conversation_users ON conversations.id = conversation_users.conversation_id").
		Where("conversations.id = ? AND conversation_users.user_id = ?", convID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *ConversationRepository) GetConversationParticipants(convID uint) ([]uint, error) {
	var userIDs []uint
	err := r.db.Table("conversation_users").
		Where("conversation_id = ?", convID).
		Pluck("user_id", &userIDs).Error
	return userIDs, err
}

func (r *ConversationRepository) GetConversationsWithUnreadCount(userID uint) ([]struct {
	Conversation *models.Conversation
	UnreadCount  int
}, error) {
	type Result struct {
		models.Conversation
		UnreadCount int
	}

	var results []Result
	err := r.db.Table("conversations").
		Select("conversations.*, COUNT(unread_messages.message_id) as unread_count").
		Joins("JOIN conversation_users ON conversations.id = conversation_users.conversation_id").
		Joins("LEFT JOIN unread_messages ON conversations.id = unread_messages.conversation_id AND unread_messages.user_id = ?", userID).
		Where("conversation_users.user_id = ?", userID).
		Group("conversations.id").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	var output []struct {
		Conversation *models.Conversation
		UnreadCount  int
	}
	for _, r := range results {
		output = append(output, struct {
			Conversation *models.Conversation
			UnreadCount  int
		}{
			Conversation: &models.Conversation{
				ID:    r.ID,
				Users: r.Users,
			},
			UnreadCount: r.UnreadCount,
		})
	}

	return output, nil
}

func (r *ConversationRepository) GetRecentConversations(userID uint, limit int) ([]models.Conversation, error) {
	var conversations []models.Conversation
	err := r.db.
		Select("conversations.id").
		Preload("Users", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username, email")
		}).
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Table("messages").
				Select("messages.*").
				Joins("JOIN (SELECT conversation_id, MAX(message_date) as max_date FROM messages GROUP BY conversation_id) latest ON messages.conversation_id = latest.conversation_id AND messages.message_date = latest.max_date").
				Order("messages.message_date DESC")
		}).
		Joins("JOIN conversation_users cu ON cu.conversation_id = conversations.id").
		Where("cu.user_id = ?", userID).
		Find(&conversations).Error

	if err != nil {
		return nil, err
	}

	sort.Slice(conversations, func(i, j int) bool {
		if len(conversations[i].Messages) == 0 {
			return false
		}
		if len(conversations[j].Messages) == 0 {
			return true
		}
		return conversations[i].Messages[0].MessageDate.After(conversations[j].Messages[0].MessageDate)
	})

	if limit > 0 && len(conversations) > limit {
		conversations = conversations[:limit]
	}

	return conversations, nil
}
