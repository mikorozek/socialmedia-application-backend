// internal/Auth/usecases/login.go
package usecases

import (
	"errors"
	"socialmedia-app/internal/shared/models"
	"socialmedia-app/internal/shared/repositories"

	"golang.org/x/crypto/bcrypt"
)

type LoginUsecase struct {
	userRepo *repositories.UserRepository
	// Jira test
}

func NewLoginUsecase() *LoginUsecase {
	return &LoginUsecase{
		userRepo: repositories.NewUserRepository(),
	}
}

func (u *LoginUsecase) Execute(email, password string) (*models.User, error) {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
