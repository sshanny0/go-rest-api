package handler

import (
	"errors"
	"net/http"
	"rest-api/internal/dto"
	"rest-api/internal/middleware"
	"rest-api/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	authService *service.AuthService
}

func NewUserHandler(authService *service.AuthService) *UserHandler {
	return &UserHandler{
		authService: authService,
	}
}

// Register Handles user registration
// @Summary Register new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.UserRegisterRequest true "User registration data"
// @Success 201 {object} dto.SuccessResponse{data=dto.UserResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/register
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	// parse  and validate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid Input",
			Error:   err.Error(),
		})
		return
	}

	// valid json, call user service to register
	user, err := h.authService.Register(req)
	if err != nil {
		// map error
		statusCode := http.StatusInternalServerError
		message := "Failed to register user"

		if errors.Is(err, service.ErrUserExists) {
			statusCode = http.StatusBadRequest
			message = err.Error()
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Message: message,
		})
		return
	}

	// tidak 500 dan req ok tidak bad request
	c.JSON(http.StatusCreated, dto.SuccessResponse{
		Success: true,
		Message: "User registered successfully",
		Data:    user,
	})
}

// Login handles user login
// @Summary User login
// @Description Login with username and password, returns JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body dto.UserLoginRequest true "Login credentials"
// @Success 200 {object} dto.SuccessResponse{data=dto.LoginResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	// Parse and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid input",
			Error:   err.Error(),
		})
		return
	}

	// Call service
	authResp, err := h.authService.Login(req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := "Failed to login"

		if errors.Is(err, service.ErrInvalidCredentials) {
			statusCode = http.StatusUnauthorized
			message = err.Error()
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Message: message,
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Login successful",
		Data:    authResp,
	})
}

// GetProfile handles get user profile (requires auth)
// @Summary Get user profile
// @Description Get authenticated user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse{data=dto.UserResponse}
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	// Get user ID from JWT (set by auth middleware)
	userID := middleware.GetUserID(c)

	// Call service
	user, err := h.authService.GetProfile(userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := "Failed to get profile"

		if errors.Is(err, service.ErrUserNotFound) {
			statusCode = http.StatusNotFound
			message = err.Error()
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Message: message,
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Profile retrieved successfully",
		Data:    user,
	})
}

// UpdateProfile handles update user profile (requires auth)
// @Summary Update user profile
// @Description Update authenticated user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profile body dto.UserUpdateRequest true "Profile update data"
// @Success 200 {object} dto.SuccessResponse{data=dto.UserResponse}
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// Get user ID from JWT
	userID := middleware.GetUserID(c)

	var req dto.UserUpdateRequest

	// Parse and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "Invalid input",
			Error:   err.Error(),
		})
		return
	}

	// Call service
	user, err := h.authService.UpdateProfile(userID, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		message := "Failed to update profile"

		if errors.Is(err, service.ErrUserNotFound) {
			statusCode = http.StatusNotFound
			message = err.Error()
		} else if errors.Is(err, service.ErrUserExists) {
			statusCode = http.StatusBadRequest
			message = err.Error()
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Message: message,
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Profile updated successfully",
		Data:    user,
	})
}

// ResetPasswordRequest initiates a reset email
// @Summary Request password reset
// @Description Send a password reset email when email exists
// @Tags auth
// @Accept json
// @Produce json
// @Param body body dto.ResetPasswordRequest true "Email"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/reset-password [post]
func (h *UserHandler) ResetPasswordRequest(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Message: "Invalid input", Error: err.Error()})
		return
	}

	cfgServerPort := "8080"
	// try to get from env via config if available
	// avoid importing config here to keep handler lightweight; use env fallback
	if port := c.Request.URL.Port(); port != "" {
		cfgServerPort = port
	}

	if err := h.authService.RequestPasswordReset(req.Email, cfgServerPort); err != nil {
		status := http.StatusInternalServerError
		msg := "Failed to request password reset"
		if errors.Is(err, service.ErrUserNotFound) {
			status = http.StatusNotFound
			msg = err.Error()
		}
		c.JSON(status, dto.ErrorResponse{Success: false, Message: msg})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Message: "If the email exists, a reset link was sent"})
}

// ResetPasswordConfirm completes the reset using token and new password
// @Summary Complete password reset
// @Description Reset password using token and new password
// @Tags auth
// @Accept json
// @Produce json
// @Param body body dto.ResetPasswordConfirmRequest true "Token and new password"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/reset-password/confirm [post]
func (h *UserHandler) ResetPasswordConfirm(c *gin.Context) {
	var req dto.ResetPasswordConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Message: "Invalid input", Error: err.Error()})
		return
	}

	if err := h.authService.ResetPassword(req.Token, req.NewPassword); err != nil {
		status := http.StatusInternalServerError
		msg := "Failed to reset password"
		if errors.Is(err, service.ErrUserNotFound) {
			status = http.StatusNotFound
			msg = err.Error()
		}
		c.JSON(status, dto.ErrorResponse{Success: false, Message: msg, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Message: "Password reset successful"})
}
