package services

import (
	"database/sql"
	"errors"

	"github.com/grahagandangr/nexboard-be/models"
	"github.com/grahagandangr/nexboard-be/repositories"
	"github.com/grahagandangr/nexboard-be/utils"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// Register creates a new user account
func (s *AuthService) Register(req *models.RegisterRequest) (*models.UserResponse, error) {
	// Check if user already exists
	_, err := s.userRepo.GetUserByEmail(req.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		ExternalID: utils.GenerateUUID(),
		Name:       req.Name,
		Email:      req.Email,
		Password:   hashedPassword,
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ExternalID: user.ExternalID,
		Name:       user.Name,
		Email:      user.Email,
	}, nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(email, password string) (*models.LoginResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// Check password
	if err := utils.CheckPassword(user.Password, password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ExternalID, user.Email)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token: token,
		User: models.UserResponse{
			ExternalID: user.ExternalID,
			Name:       user.Name,
			Email:      user.Email,
		},
	}, nil
}

// GetProfile retrieves user profile using their external ID from token claims
func (s *AuthService) GetProfile(externalID string) (*models.User, error) {
	user, err := s.userRepo.GetUserByExternalID(externalID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	// Sanitize output
	user.Password = ""
	user.ID = 0
	return user, nil
}

// UpdateProfile updates authenticated user's profile
func (s *AuthService) UpdateProfile(externalID string, req *models.UpdateProfileRequest) (*models.User, error) {
	user, err := s.userRepo.GetUserByExternalID(externalID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	user.Name = req.Name
	user.AvatarURL = req.AvatarURL

	if err := s.userRepo.UpdateUserProfile(user); err != nil {
		return nil, err
	}

	// Sanitize
	user.Password = ""
	user.ID = 0
	return user, nil
}
