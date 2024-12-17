package repositories

import (
	"errors"
	"socialmedia-backend/internal/shared/db"
	"socialmedia-backend/internal/shared/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: db.GetDB(),
	}
}

// GetByEmail 
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) Create(user *models.User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		return errors.New("failed to create user")
	}
	return nil
}
