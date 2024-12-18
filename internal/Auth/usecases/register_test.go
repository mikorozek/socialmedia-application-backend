package usecases

import (
	"testing"
)

func TestRegisterUsecase_Execute(t *testing.T) {
	mockRepo := NewMockUserRepository()
	registerUsecase := &RegisterUsecase{
		userRepo: mockRepo,
	}

	t.Run("Successful Registration", func(t *testing.T) {
		err := registerUsecase.Execute("newuser", "new@example.com", "password123")
		if err != nil {
			t.Errorf("Expected successful registration, got error: %v", err)
		}

		user, err := mockRepo.GetByEmail("new@example.com")
		if err != nil {
			t.Error("Expected to find created user, got error")
		}
		if user.Username != "newuser" {
			t.Errorf("Expected username %s, got %s", "newuser", user.Username)
		}
		if user.Email != "new@example.com" {
			t.Errorf("Expected email %s, got %s", "new@example.com", user.Email)
		}
	})

	t.Run("Email Already Exists", func(t *testing.T) {

		err := registerUsecase.Execute("firstuser", "test@example.com", "password123")
		if err != nil {
			t.Fatalf("Failed to create first user: %v", err)
		}

		err = registerUsecase.Execute("seconduser", "test@example.com", "password456")
		if err == nil {
			t.Error("Expected error for duplicate email, got nil")
		}
		if err.Error() != "email already registered" {
			t.Errorf("Expected 'email already registered' error, got: %v", err)
		}
	})

	t.Run("Password Hashing", func(t *testing.T) {
		password := "testpassword"
		err := registerUsecase.Execute("hashtest", "hash@example.com", password)
		if err != nil {
			t.Fatalf("Failed to register user: %v", err)
		}

		user, _ := mockRepo.GetByEmail("hash@example.com")
		if user.PasswordHash == password {
			t.Error("Password was not hashed")
		}
		if len(user.PasswordHash) == 0 {
			t.Error("Password hash is empty")
		}
	})
}
