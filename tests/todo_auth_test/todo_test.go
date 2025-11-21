package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
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

// TodoTestSuite defines the test suite
type TodoTestSuite struct {
	suite.Suite
	db     *gorm.DB
	router *gin.Engine
	token  string
	userID uint
}

// SetupSuite runs once before all tests
func (suite *TodoTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	cfg := config.LoadConfig()

	dsn := "host=" + cfg.DBHost +
		" user=" + cfg.DBUser +
		" password=" + cfg.DBPassword +
		" dbname=" + cfg.DBName +
		" port=" + cfg.DBPort +
		" sslmode=disable TimeZone=Asia/Jakarta"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	suite.Require().NoError(err)

	suite.db = db

	err = db.AutoMigrate(&model.User{}, &model.Todo{})
	suite.Require().NoError(err)

	// Initialize dependencies
	userRepo := repository.NewUserRepository(db)
	todoRepo := repository.NewTodoRepository(db)
	authService := service.NewAuthService(userRepo)
	todoService := service.NewTodoService(todoRepo)
	userHandler := handler.NewUserHandler(authService)
	todoHandler := handler.NewTodoHandler(todoService)
	healthHandler := handler.NewHealthHandler(db)

	router := gin.New()
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.CORSMiddleware())
	route.SetupRoutes(router, userHandler, healthHandler, todoHandler)

	suite.router = router

	// Create test user and get token
	suite.createTestUserAndToken()
}

// createTestUserAndToken creates a test user and stores token
func (suite *TodoTestSuite) createTestUserAndToken() {
	// Register
	regBody := dto.RegisterRequest{
		Username: "todotest",
		Email:    "todotest@example.com",
		Password: "password123",
		FullName: "Todo Test User",
	}
	jsonReg, _ := json.Marshal(regBody)
	reqReg := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(jsonReg))
	reqReg.Header.Set("Content-Type", "application/json")
	wReg := httptest.NewRecorder()
	suite.router.ServeHTTP(wReg, reqReg)

	var regResponse dto.SuccessResponse
	json.Unmarshal(wReg.Body.Bytes(), &regResponse)
	userData := regResponse.Data.(map[string]interface{})
	suite.userID = uint(userData["id"].(float64))

	// Login
	loginBody := dto.LoginRequest{
		Username: "todotest",
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
	suite.token = loginData["token"].(string)
}

// TearDownSuite runs once after all tests
func (suite *TodoTestSuite) TearDownSuite() {
	suite.db.Exec("DELETE FROM todos")
	suite.db.Exec("DELETE FROM users")
}

// SetupTest runs before each test
func (suite *TodoTestSuite) SetupTest() {
	suite.db.Exec("DELETE FROM todos WHERE user_id = ?", suite.userID)
}

// TestCreateTodo tests creating a new todo
func (suite *TodoTestSuite) TestCreateTodo() {
	reqBody := dto.CreateTodoRequest{
		Title:       "Test Todo",
		Description: "Test description",
		Status:      "pending",
		Priority:    "high",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response dto.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)

	todoData, ok := response.Data.(map[string]interface{})
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), "Test Todo", todoData["title"])
	assert.Equal(suite.T(), "pending", todoData["status"])
}

// TestGetAllTodos tests getting all user's todos
func (suite *TodoTestSuite) TestGetAllTodos() {
	// Create test todos
	suite.createTestTodo("Todo 1", "pending", "high")
	suite.createTestTodo("Todo 2", "in_progress", "medium")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/todos", nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response dto.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)

	todos, ok := response.Data.([]interface{})
	assert.True(suite.T(), ok)
	assert.GreaterOrEqual(suite.T(), len(todos), 2)
}

// TestGetTodoByID tests getting a specific todo
func (suite *TodoTestSuite) TestGetTodoByID() {
	todoID := suite.createTestTodo("Test Todo", "pending", "high")

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/todos/%d", todoID), nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response dto.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)

	todoData, ok := response.Data.(map[string]interface{})
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), "Test Todo", todoData["title"])
}

// TestUpdateTodo tests updating a todo
func (suite *TodoTestSuite) TestUpdateTodo() {
	todoID := suite.createTestTodo("Original Title", "pending", "low")

	updateBody := map[string]interface{}{
		"title":  "Updated Title",
		"status": "completed",
	}
	jsonBody, _ := json.Marshal(updateBody)

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/todos/%d", todoID), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response dto.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)

	todoData, ok := response.Data.(map[string]interface{})
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), "Updated Title", todoData["title"])
	assert.Equal(suite.T(), "completed", todoData["status"])
}

// TestDeleteTodo tests deleting a todo
func (suite *TodoTestSuite) TestDeleteTodo() {
	todoID := suite.createTestTodo("To Delete", "pending", "low")

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/todos/%d", todoID), nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response dto.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), response.Success)

	// Verify todo is deleted (soft delete)
	var todo model.Todo
	err = suite.db.Where("id = ?", todoID).First(&todo).Error
	assert.Error(suite.T(), err) // Should not find (soft deleted)
}

// Helper function to create test todo
func (suite *TodoTestSuite) createTestTodo(title, status, priority string) uint {
	reqBody := dto.CreateTodoRequest{
		Title:    title,
		Status:   status,
		Priority: priority,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var response dto.SuccessResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	todoData := response.Data.(map[string]interface{})
	return uint(todoData["id"].(float64))
}

// Run the test suite
func TestTodoTestSuite(t *testing.T) {
	suite.Run(t, new(TodoTestSuite))
}
