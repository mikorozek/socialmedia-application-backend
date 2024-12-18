package usecases

import (
	"errors"
	"socialmedia-backend/internal/shared/models"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	users map[string]*models.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*models.User),
	}
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	user, exists := m.users[email]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *MockUserRepository) Create(user *models.User) error {
	m.users[user.Email] = user
	return nil
}

func TestLoginUsecase_Execute(t *testing.T) {

	mockRepo := NewMockUserRepository()
	password := "testpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	testUser := &models.User{
		ID:           1,
		Email:        "test@example.com",
		Username:     "testuser",
		PasswordHash: string(hashedPassword),
	}

	mockRepo.Create(testUser)

	loginUsecase := &LoginUsecase{
		userRepo: mockRepo,
	}

	t.Run("Successful Login", func(t *testing.T) {
		user, err := loginUsecase.Execute("test@example.com", password)
		if err != nil {
			t.Errorf("Expected successful login, got error: %v", err)
		}
		if user == nil {
			t.Error("Expected user object, got nil")
		}
		if user.Email != testUser.Email {
			t.Errorf("Expected email %s, got %s", testUser.Email, user.Email)
		}
	})

	t.Run("Invalid Password", func(t *testing.T) {
		_, err := loginUsecase.Execute("test@example.com", "wrongpassword")
		if err == nil {
			t.Error("Expected error for invalid password, got nil")
		}
	})

	t.Run("Non-existent User", func(t *testing.T) {
		_, err := loginUsecase.Execute("nonexistent@example.com", password)
		if err == nil {
			t.Error("Expected error for non-existent user, got nil")
		}
	})
}
