package repository

import (
	"errors"
	"rest-api/internal/model"
	"time"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

// inisiasi repo
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// create to db
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// find by ID
func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByUsername retrieves user by username
func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil user, no error for not found
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail retrieves user by email
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil user, no error for not found
		}
		return nil, err
	}
	return &user, nil
}

// Update updates user data
func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Delete soft deletes a user
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

// ExistsByUsername checks if username already exists
func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// ExistsByEmail checks if email already exists
func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// FindByResetToken finds a user by reset token
func (r *UserRepository) FindByResetToken(token string) (*model.User, error) {
	var user model.User
	err := r.db.Where("reset_password_token = ?", token).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// SaveResetToken saves reset token and expiry for a user
func (r *UserRepository) SaveResetToken(userID uint, token string, expiry *time.Time) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"reset_password_token":  token,
		"reset_password_expiry": expiry,
	}).Error
}
