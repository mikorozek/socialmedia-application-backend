package usecases

import (
	"errors"
	"socialmedia-backend/internal/shared/models"
	"socialmedia-backend/internal/shared/repositories"
	"time"
)

type ConversationUsecase struct {
	convRepo *repositories.ConversationRepository
	userRepo *repositories.UserRepository
}

func NewConversationUsecase() *ConversationUsecase {
	return &ConversationUsecase{
		convRepo: repositories.NewConversationRepository(),
		userRepo: repositories.NewUserRepository(),
	}
}

func (u *ConversationUsecase) CreateConversation(userIDs []uint) (*models.Conversation, error) {
	if len(userIDs) < 2 {
		return nil, errors.New("conversation requires at least 2 users")
	}

	conv := &models.Conversation{}
	err := u.convRepo.Create(conv)
	if err != nil {
		return nil, errors.New("failed to create conversation")
	}

	for _, userID := range userIDs {
		if err := u.convRepo.AddUserToConversation(conv.ID, userID); err != nil {
			return nil, err
		}
	}

	return conv, nil
}

func (u *ConversationUsecase) SendMessage(conversationID uint, senderID uint, content string, photoURL string) error {
	if err := u.checkUserAccess(conversationID, senderID); err != nil {
		return err
	}

	message := &models.Message{
		ConversationID: conversationID,
		UserID:         senderID,
		Content:        content,
		PhotoURL:       photoURL,
		MessageDate:    time.Now(),
	}

	if err := u.convRepo.AddMessage(message); err != nil {
		return err
	}

	// Dodaj wiadomość jako nieprzeczytaną dla wszystkich innych uczestników konwersacji
	participants, err := u.convRepo.GetConversationParticipants(conversationID)
	if err != nil {
		return err
	}

	for _, participantID := range participants {
		if participantID != senderID {
			if err := u.convRepo.AddUnreadMessage(message.ID, participantID, conversationID); err != nil {
				return err
			}
		}
	}

	return nil
}

func (u *ConversationUsecase) GetConversationMessages(conversationID uint, userID uint, limit int, offset int) ([]models.Message, error) {
	if err := u.checkUserAccess(conversationID, userID); err != nil {
		return nil, err
	}

	messages, err := u.convRepo.GetMessages(conversationID, limit, offset)
	if err != nil {
		return nil, errors.New("failed to get messages")
	}

	return messages, nil
}

func (u *ConversationUsecase) GetUserConversations(userID uint) ([]models.Conversation, error) {
	conversations, err := u.convRepo.GetUserConversations(userID)
	if err != nil {
		return nil, errors.New("failed to get user conversations")
	}

	return conversations, nil
}

func (u *ConversationUsecase) EditMessage(messageID uint, userID uint, newContent string) error {
	message, err := u.convRepo.GetMessageByID(messageID)
	if err != nil {
		return errors.New("message not found")
	}

	if message.UserID != userID {
		return errors.New("cannot edit message from another user")
	}

	message.Content = newContent
	return u.convRepo.UpdateMessage(message)
}

func (u *ConversationUsecase) DeleteMessage(messageID uint, userID uint) error {
	message, err := u.convRepo.GetMessageByID(messageID)
	if err != nil {
		return errors.New("message not found")
	}

	if message.UserID != userID {
		return errors.New("cannot delete message from another user")
	}

	if err := u.convRepo.DeleteUnreadMessages(messageID); err != nil {
		return err
	}

	return u.convRepo.DeleteMessage(messageID)
}

func (u *ConversationUsecase) checkUserAccess(conversationID uint, userID uint) error {
	exists, err := u.convRepo.CheckUserInConversation(conversationID, userID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("user does not have access to this conversation")
	}
	return nil
}

func (u *ConversationUsecase) MarkMessageAsRead(messageID uint, userID uint) error {
	message, err := u.convRepo.GetMessageByID(messageID)
	if err != nil {
		return errors.New("message not found")
	}

	return u.convRepo.DeleteUnreadMessage(messageID, userID, message.ConversationID)
}

func (u *ConversationUsecase) GetUnreadMessages(userID uint) ([]models.Message, error) {
	return u.convRepo.GetUnreadMessagesForUser(userID)
}

func (u *ConversationUsecase) GetConversationsWithUnreadCount(userID uint) ([]struct {
	Conversation *models.Conversation
	UnreadCount  int
}, error) {
	return u.convRepo.GetConversationsWithUnreadCount(userID)
}
