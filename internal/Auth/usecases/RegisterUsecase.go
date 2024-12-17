package usecases

import (
	"errors"
	"socialmedia-backend/internal/shared/models"
	"socialmedia-backend/internal/shared/repositories"

	"golang.org/x/crypto/bcrypt"
)

type RegisterUsecase struct {
	userRepo *repositories.UserRepository
}

func NewRegisterUsecase() *RegisterUsecase {
	return &RegisterUsecase{
		userRepo: repositories.NewUserRepository(),
	}
}

func (u *RegisterUsecase) Execute(username, email, password string) error {
	if _, err := u.userRepo.GetByEmail(email); err == nil {
		return errors.New("email already registered")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	newUser := &models.User{
		Username:       username,
		Email:          email,
		PasswordHash: 	string(passwordHash),
	}

	if err := u.userRepo.Create(newUser); err != nil {
		return errors.New("failed to create user")
	}

	return nil
}
