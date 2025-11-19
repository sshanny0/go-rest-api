package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"rest-api/internal/dto"
	"rest-api/internal/model"
	"rest-api/internal/repository"
	"rest-api/internal/utils"
	"time"
)

var (
	ErrUserExists         = errors.New("username already exists")
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserNotFound       = errors.New("user not found")
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// Register User
func (service *AuthService) Register(req dto.RegisterRequest) (*dto.UserResponse, error) {
	// 1. Check if username already exist
	existsUsername, err := service.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username: %w", err)
	}
	if existsUsername {
		return nil, ErrUserExists
	}

	// Business Rule 2: Check if email already exist
	existsEmail, err := service.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if existsEmail {
		return nil, ErrEmailExists
	}

	// Business rule 3: New Username, New Email (Valid New User)
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// create user
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Fullname: req.FullName,
	}

	// save to database via repository
	if err := service.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// convert to DTO response
	return service.toUserResponse(user), nil
}

// Login melakukan autentikasi user dan mengembalikan JWT token
func (s *AuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	// 1. Find user by username
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// 2. Verify password
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	// 3.Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Return auth response with token and user data
	return &dto.LoginResponse{
		Token: token,
		User:  *s.toUserResponse(user),
	}, nil
}

// GetProfile mendapatkan profile user berdasarkan ID
func (s *AuthService) GetProfile(userID uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return s.toUserResponse(user), nil
}

// UpdateProfile mengupdate profile user
func (s *AuthService) UpdateProfile(userID uint, req dto.UserUpdateRequest) (*dto.UserResponse, error) {
	// Find existing user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// Business Rule: Check if email already used by another user
	if req.Email != "" && req.Email != user.Email {
		existingUser, err := s.userRepo.FindByEmail(req.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to check email: %w", err)
		}
		if existingUser != nil && existingUser.ID != userID {
			return nil, ErrUserExists
		}
		user.Email = req.Email
	}

	// Update fields
	if req.FullName != "" {
		user.Fullname = req.FullName
	}

	// Save changes
	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return s.toUserResponse(user), nil
}

func (s *AuthService) toUserResponse(user *model.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.Fullname,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// RequestPasswordReset initiates a password reset flow for the given email.
// It generates a secure token, saves it with expiry, and sends a reset email.
func (s *AuthService) RequestPasswordReset(email string, serverPort string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return fmt.Errorf("failed to lookup user: %w", err)
	}
	if user == nil {
		// Do not reveal existence in production; following prompt, return not found
		return ErrUserNotFound
	}

	// generate secure token
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}
	token := base64.RawURLEncoding.EncodeToString(b)

	// expiry (1 hour)
	expiry := time.Now().Add(1 * time.Hour)

	// save token and expiry
	if err := s.userRepo.SaveResetToken(user.ID, token, &expiry); err != nil {
		return fmt.Errorf("failed to save reset token: %w", err)
	}

	// build reset link (API confirm endpoint expects POST; this link contains token)
	resetLink := fmt.Sprintf("http://localhost:%s/api/v1/auth/reset-password/confirm?token=%s", serverPort, token)

	// send email (placeholder)
	if err := utils.SendResetEmail(user.Email, resetLink); err != nil {
		return fmt.Errorf("failed to send reset email: %w", err)
	}

	return nil
}

// ResetPassword completes password reset using token and new password
func (s *AuthService) ResetPassword(token, newPassword string) error {
	user, err := s.userRepo.FindByResetToken(token)
	if err != nil {
		return fmt.Errorf("failed to lookup token: %w", err)
	}
	if user == nil {
		return ErrUserNotFound
	}
	if user.ResetPasswordExpiry == nil || time.Now().After(*user.ResetPasswordExpiry) {
		return errors.New("reset token expired or invalid")
	}

	// hash new password
	hashed, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = hashed
	user.ResetPasswordToken = ""
	user.ResetPasswordExpiry = nil

	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update user password: %w", err)
	}

	return nil
}
