package services

import (
	"errors"

	"job-portal/models"
	"job-portal/repository"
	"job-portal/utils"

	"gorm.io/gorm"
)

// RegisterInput holds the data required for user registration
type RegisterInput struct {
	Name     string      `json:"name" binding:"required,min=2,max=100"`
	Email    string      `json:"email" binding:"required,email"`
	Password string      `json:"password" binding:"required,min=6"`
	Role     models.Role `json:"role" binding:"required,oneof=ADMIN RECRUITER CANDIDATE"`
}

// LoginInput holds the data required for user login
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse is returned after successful login/register
type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// AuthService defines authentication business logic
type AuthService interface {
	Register(input RegisterInput) (*AuthResponse, error)
	Login(input LoginInput) (*AuthResponse, error)
}

type authService struct {
	userRepo repository.UserRepository
}

// NewAuthService creates a new AuthService instance
func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(input RegisterInput) (*AuthResponse, error) {
	// Check if email already exists
	existing, err := s.userRepo.FindByEmail(input.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("database error")
	}
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashed, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashed,
		Role:     input.Role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	token, err := utils.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &AuthResponse{Token: token, User: *user}, nil
}

func (s *authService) Login(input LoginInput) (*AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &AuthResponse{Token: token, User: *user}, nil
}
