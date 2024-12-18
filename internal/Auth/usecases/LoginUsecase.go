// internal/Auth/usecases/login.go
package usecases

import (
	"errors"
	"socialmedia-backend/internal/shared/models"
	"socialmedia-backend/internal/shared/repositories"

	"golang.org/x/crypto/bcrypt"
)

type LoginUsecase struct {
	userRepo repositories.UserRepositoryInterface
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

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
