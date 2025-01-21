package usecases

import (
	"errors"
	"socialmedia-backend/internal/shared/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserProfileUsecase struct {
	userRepo *repositories.UserRepository
}

type UserProfileResponse struct {
	ID          uint   `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Description string `json:"description"`
}

func NewUserProfileUsecase() *UserProfileUsecase {
	return &UserProfileUsecase{
		userRepo: repositories.NewUserRepository(),
	}
}

func (u *UserProfileUsecase) GetUserProfile(userID uint) (*UserProfileResponse, error) {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return &UserProfileResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Description: user.Description,
	}, nil
}

func (u *UserProfileUsecase) UpdateUserProfile(userID uint, username, password, description string, requestBody map[string]interface{}) error {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	if _, exists := requestBody["username"]; exists {
		user.Username = username
	}

	if _, exists := requestBody["description"]; exists {
		user.Description = description
	}

	if _, exists := requestBody["password"]; exists && password != "" {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("failed to hash password")
		}
		user.PasswordHash = string(passwordHash)
	}

	return u.userRepo.Update(user)
}
