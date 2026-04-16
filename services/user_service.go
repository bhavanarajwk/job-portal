package services

import (
	"errors"

	"job-portal/models"
	"job-portal/repository"
	"job-portal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserService defines user management business logic (admin)
type UserService interface {
	GetAllUsers(pagination utils.PaginationParams) (utils.PaginatedResponse, error)
	DeleteUser(id uuid.UUID) error
	GetMe(userID uuid.UUID) (*models.User, error)
	UpdateMe(userID uuid.UUID, input UpdateMeInput) (*models.User, error)
	ChangePassword(userID uuid.UUID, input ChangePasswordInput) error
}

type UpdateMeInput struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=100"`
	Email       string `json:"email" binding:"omitempty,email"`
	Password    string `json:"password" binding:"omitempty,min=6"`
	NewPassword string `json:"new_password" binding:"omitempty,min=6"`
	NewPasswordAlt string `json:"newPassword" binding:"omitempty,min=6"`
}

type ChangePasswordInput struct {
	CurrentPassword string `json:"current_password" binding:"required,min=6"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
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

func (s *userService) GetMe(userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("database error")
	}

	return user, nil
}

func (s *userService) UpdateMe(userID uuid.UUID, input UpdateMeInput) (*models.User, error) {
	passwordToSet := input.NewPassword
	if passwordToSet == "" {
		passwordToSet = input.NewPasswordAlt
	}
	if passwordToSet == "" {
		passwordToSet = input.Password
	}

	if input.Name == "" && input.Email == "" && passwordToSet == "" {
		return nil, errors.New("at least one field is required")
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("database error")
	}

	if input.Email != "" && input.Email != user.Email {
		existing, err := s.userRepo.FindByEmail(input.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("database error")
		}
		if existing != nil && existing.ID != userID {
			return nil, errors.New("email already registered")
		}
		user.Email = input.Email
	}

	if input.Name != "" {
		user.Name = input.Name
	}

	if passwordToSet != "" {
		hashed, err := utils.HashPassword(passwordToSet)
		if err != nil {
			return nil, errors.New("failed to hash password")
		}
		user.Password = hashed
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.New("failed to update user")
	}

	return user, nil
}

func (s *userService) ChangePassword(userID uuid.UUID, input ChangePasswordInput) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return errors.New("database error")
	}

	if !utils.CheckPasswordHash(input.CurrentPassword, user.Password) {
		return errors.New("invalid current password")
	}

	if utils.CheckPasswordHash(input.NewPassword, user.Password) {
		return errors.New("new password must be different from current password")
	}

	hashed, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user.Password = hashed

	if err := s.userRepo.Update(user); err != nil {
		return errors.New("failed to update password")
	}

	return nil
}
