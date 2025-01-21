package usecases

import (
	"errors"
	"socialmedia-backend/internal/shared/repositories"
)

type UserSearchUsecase struct {
	userRepo *repositories.UserRepository
}

func NewUserSearchUsecase() *UserSearchUsecase {
	return &UserSearchUsecase{
		userRepo: repositories.NewUserRepository(),
	}
}

type UserSearchResult struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (u *UserSearchUsecase) SearchUsers(query string) ([]UserSearchResult, error) {
	if len(query) < 3 {
		return nil, errors.New("search query must be at least 3 characters long")
	}

	users, err := u.userRepo.SearchUsers(query)
	if err != nil {
		return nil, err
	}

	var results []UserSearchResult
	for _, user := range users {
		results = append(results, UserSearchResult{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		})
	}

	return results, nil
}
