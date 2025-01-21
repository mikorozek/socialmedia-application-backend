package usecases

import (
	"errors"
	"reflect"
	"socialmedia-backend/internal/shared/models"
	"socialmedia-backend/internal/shared/repositories"
	"socialmedia-backend/internal/shared/services"
	"sort"
	"time"
)

type ConversationUsecase struct {
	convRepo   *repositories.ConversationRepository
	msgRepo    *repositories.MessageRepository
	unreadRepo *repositories.UnreadConversationRepository
	userRepo   *repositories.UserRepository
	wsService  *services.WebSocketService
}

func NewConversationUsecase(wsService *services.WebSocketService) *ConversationUsecase {
	return &ConversationUsecase{
		convRepo:   repositories.NewConversationRepository(),
		msgRepo:    repositories.NewMessageRepository(),
		unreadRepo: repositories.NewUnreadConversationRepository(),
		userRepo:   repositories.NewUserRepository(),
		wsService:  wsService,
	}
}

func (u *ConversationUsecase) CreateConversation(userIDs []uint) (*models.Conversation, error) {
	sort.Slice(userIDs, func(i, j int) bool {
		return userIDs[i] < userIDs[j]
	})

	conversations, err := u.convRepo.GetUserConversations(userIDs[0])
	if err != nil {
		return nil, err
	}

	for _, conv := range conversations {
		if len(conv.Users) == len(userIDs) {
			existingUserIDs := make([]uint, len(conv.Users))
			for i, user := range conv.Users {
				existingUserIDs[i] = user.ID
			}
			sort.Slice(existingUserIDs, func(i, j int) bool {
				return existingUserIDs[i] < existingUserIDs[j]
			})

			if reflect.DeepEqual(existingUserIDs, userIDs) {
				return nil, errors.New("conversation between these users already exists")
			}
		}
	}
	if len(userIDs) < 2 {
		return nil, errors.New("conversation requires at least 2 users")
	}

	conv := &models.Conversation{}
	err = u.convRepo.Create(conv)
	if err != nil {
		return nil, errors.New("failed to create conversation")
	}

	for _, userID := range userIDs {
		if err := u.convRepo.AddUserToConversation(conv.ID, userID); err != nil {
			return nil, err
		}
	}

	fullConversation, err := u.convRepo.GetByID(conv.ID)
	if err != nil {
		return nil, err
	}

	return fullConversation, nil
}

func (u *ConversationUsecase) SendMessage(conversationID uint, senderID uint, content string, photoURL string) (*models.Message, error) {
	if err := u.checkUserAccess(conversationID, senderID); err != nil {
		return nil, err
	}

	message := &models.Message{
		ConversationID: conversationID,
		UserID:         senderID,
		Content:        content,
		PhotoURL:       photoURL,
		MessageDate:    time.Now(),
	}

	if err := u.msgRepo.AddMessage(message); err != nil {
		return nil, err
	}

	sender, err := u.userRepo.GetByID(senderID)
	if err != nil {
		return nil, err
	}

	participants, err := u.convRepo.GetConversationParticipants(conversationID)
	if err != nil {
		return nil, err
	}

	var recipients []uint
	for _, participantID := range participants {
		if participantID != senderID {
			recipients = append(recipients, participantID)

			if err := u.unreadRepo.UpdateUnreadConversation(
				conversationID,
				participantID,
				content,
				message.MessageDate,
			); err != nil {
				return nil, err
			}
		}
	}

	notification := services.MessageNotification{
		ConversationID: conversationID,
		SenderID:       senderID,
		SenderUsername: sender.Username,
		Content:        content,
	}

	u.wsService.NotifyUsers(notification, recipients)

	return message, nil
}

func (u *ConversationUsecase) GetConversationMessages(conversationID uint, userID uint, limit int, offset int) ([]models.Message, error) {
	if err := u.checkUserAccess(conversationID, userID); err != nil {
		return nil, err
	}

	if err := u.unreadRepo.MarkAsRead(conversationID, userID); err != nil {
		return nil, err
	}

	messages, err := u.msgRepo.GetMessages(conversationID, limit, offset)
	if err != nil {
		return nil, errors.New("failed to get messages")
	}

	return messages, nil
}

func (u *ConversationUsecase) EditMessage(messageID uint, userID uint, newContent string) error {
	message, err := u.msgRepo.GetMessageByID(messageID)
	if err != nil {
		return errors.New("message not found")
	}

	if message.UserID != userID {
		return errors.New("cannot edit message from another user")
	}

	message.Content = newContent
	return u.msgRepo.UpdateMessage(message)
}

func (u *ConversationUsecase) DeleteMessage(messageID uint, userID uint) error {
	message, err := u.msgRepo.GetMessageByID(messageID)
	if err != nil {
		return errors.New("message not found")
	}

	if message.UserID != userID {
		return errors.New("cannot delete message from another user")
	}

	return u.msgRepo.DeleteMessage(messageID)
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

func (u *ConversationUsecase) MarkConversationAsRead(conversationID uint, userID uint) error {
	if err := u.checkUserAccess(conversationID, userID); err != nil {
		return err
	}
	return u.unreadRepo.MarkAsRead(conversationID, userID)
}

func (u *ConversationUsecase) GetUnreadConversations(userID uint) ([]models.UnreadConversation, error) {
	return u.unreadRepo.GetUnreadConversations(userID)
}

func (u *ConversationUsecase) GetRecentConversations(userID uint, limit int) ([]models.Conversation, error) {
	return u.convRepo.GetRecentConversations(userID, limit)
}
