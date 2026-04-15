package services

import (
	"errors"

	"job-portal/repository"
	"job-portal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserService defines user management business logic (admin)
type UserService interface {
	GetAllUsers(pagination utils.PaginationParams) (utils.PaginatedResponse, error)
	DeleteUser(id uuid.UUID) error
}

type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new UserService instance
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetAllUsers(pagination utils.PaginationParams) (utils.PaginatedResponse, error) {
	users, total, err := s.userRepo.FindAll(pagination.Offset, pagination.PageSize)
	if err != nil {
		return utils.PaginatedResponse{}, errors.New("failed to fetch users")
	}
	return utils.BuildPaginatedResponse(users, total, pagination), nil
}

func (s *userService) DeleteUser(id uuid.UUID) error {
	_, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return errors.New("database error")
	}

	return s.userRepo.Delete(id)
}
