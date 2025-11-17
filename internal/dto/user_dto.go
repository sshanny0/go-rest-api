package dto

import "time"

// DTO REQUEST
// Register
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email,max=100"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"fullname" binding:"required,max=100"`
}

// Login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Reset password request (initiate)
type ResetPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// Reset password confirm (complete)
type ResetPasswordConfirmRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// Update Profile
type UserUpdateRequest struct {
	Email    string `json:"email" binding:"omitempty,email,max=100"`
	FullName string `json:"full_name" binding:"omitempty,max=100"`
}

// DTO RESPONSE
// User Reponse
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"string"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Login Reponse
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// Success Reponse
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Error Reponse
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}
