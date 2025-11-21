package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"rest-api/internal/config"
	"rest-api/internal/dto"
	"rest-api/internal/handler"
	"rest-api/internal/middleware"
	"rest-api/internal/model"
	"rest-api/internal/repository"
	"rest-api/internal/route"
	"rest-api/internal/service"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// AuthTestSuite defines the test suite
type AuthTestSuite struct {
	suite.Suite
	db     *gorm.DB
	router *gin.Engine
}

// SetupSuite runs once before all tests
func (suite *AuthTestSuite) SetupSuite() {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Load configuration
	cfg := config.LoadConfig()

	// Connect to test database
	dsn := "host=" + cfg.DBHost +
		" user=" + cfg.DBUser +
		" password=" + cfg.DBPassword +
		" dbname=" + cfg.DBName +
		" port=" + cfg.DBPort +
		" sslmode=disable TimeZone=Asia/Jakarta"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	suite.Require().NoError(err, "Failed to connect to test database")

	suite.db = db

	// Auto-migrate models
	err = db.AutoMigrate(&model.User{}, &model.Todo{})
	suite.Require().NoError(err, "Failed to migrate test database")

	// Initialize dependencies
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	userHandler := handler.NewUserHandler(authService)
	healthHandler := handler.NewHealthHandler(db)

	// Dummy handlers for routes that won't be tested
	todoRepo := repository.NewTodoRepository(db)
	todoService := service.NewTodoService(todoRepo)
	todoHandler := handler.NewTodoHandler(todoService)

	// Setup router
	router := gin.New()
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware())
	route.SetupRoutes(router, userHandler, healthHandler, todoHandler)

	suite.router = router
}

// TearDownSuite runs once after all tests
func (suite *AuthTestSuite) TearDownSuite() {
	// Clean up test data
	suite.db.Exec("DELETE FROM todos")
	suite.db.Exec("DELETE FROM users")
}

// SetupTest runs before each test
func (suite *AuthTestSuite) SetupTest() {
	// Clean tables before each test
	suite.db.Exec("DELETE FROM todos")
	suite.db.Exec("DELETE FROM users")
}

// TestRegisterSuccess tests successful user registration
func (suite *AuthTestSuite) TestRegisterSuccess() {
	// Prepare request
	reqBody := dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Create HTTP request
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Record response
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response dto.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)
	assert.Contains(suite.T(), response.Message, "berhasil")

	// Verify user data in response
	userData, ok := response.Data.(map[string]interface{})
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), "testuser", userData["username"])
	assert.Equal(suite.T(), "test@example.com", userData["email"])
}

// TestRegisterDuplicateUsername tests registration with existing username
func (suite *AuthTestSuite) TestRegisterDuplicateUsername() {
	// Create first user
	reqBody := dto.RegisterRequest{
		Username: "testuser",
		Email:    "test1@example.com",
		Password: "password123",
		FullName: "Test User",
	}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	// Try to create duplicate username
	reqBody2 := dto.RegisterRequest{
		Username: "testuser", // Same username
		Email:    "test2@example.com",
		Password: "password123",
		FullName: "Test User 2",
	}
	jsonBody2, _ := json.Marshal(reqBody2)
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(jsonBody2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	suite.router.ServeHTTP(w2, req2)

	// Assert error response
	assert.Equal(suite.T(), http.StatusConflict, w2.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w2.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), response.Success)
}

// TestLoginSuccess tests successful login
func (suite *AuthTestSuite) TestLoginSuccess() {
	// Register user first
	regBody := dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
	}
	jsonReg, _ := json.Marshal(regBody)
	reqReg := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(jsonReg))
	reqReg.Header.Set("Content-Type", "application/json")
	wReg := httptest.NewRecorder()
	suite.router.ServeHTTP(wReg, reqReg)

	// Login
	loginBody := dto.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}
	jsonLogin, _ := json.Marshal(loginBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(jsonLogin))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response dto.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)

	// Verify token in response
	loginData, ok := response.Data.(map[string]interface{})
	assert.True(suite.T(), ok)
	assert.NotEmpty(suite.T(), loginData["token"])
	assert.NotNil(suite.T(), loginData["user"])
}

// TestLoginInvalidCredentials tests login with wrong password
func (suite *AuthTestSuite) TestLoginInvalidCredentials() {
	// Register user first
	regBody := dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
	}
	jsonReg, _ := json.Marshal(regBody)
	reqReg := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(jsonReg))
	reqReg.Header.Set("Content-Type", "application/json")
	wReg := httptest.NewRecorder()
	suite.router.ServeHTTP(wReg, reqReg)

	// Login with wrong password
	loginBody := dto.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}
	jsonLogin, _ := json.Marshal(loginBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(jsonLogin))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assert error response
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), response.Success)
}

// TestGetProfile tests getting user profile with authentication
func (suite *AuthTestSuite) TestGetProfile() {
	// Register and login to get token
	regBody := dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
	}
	jsonReg, _ := json.Marshal(regBody)
	reqReg := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(jsonReg))
	reqReg.Header.Set("Content-Type", "application/json")
	wReg := httptest.NewRecorder()
	suite.router.ServeHTTP(wReg, reqReg)

	loginBody := dto.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}
	jsonLogin, _ := json.Marshal(loginBody)
	reqLogin := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(jsonLogin))
	reqLogin.Header.Set("Content-Type", "application/json")
	wLogin := httptest.NewRecorder()
	suite.router.ServeHTTP(wLogin, reqLogin)

	var loginResponse dto.SuccessResponse
	json.Unmarshal(wLogin.Body.Bytes(), &loginResponse)
	loginData := loginResponse.Data.(map[string]interface{})
	token := loginData["token"].(string)

	// Get profile with token
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response dto.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)

	// Verify user data
	userData, ok := response.Data.(map[string]interface{})
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), "testuser", userData["username"])
}

// Run the test suite
func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
