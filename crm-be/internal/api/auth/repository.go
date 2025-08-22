package authapi

import (
	"skripsi-be/internal/models/entities"

	"gorm.io/gorm"
)

type AuthRepositoryInterface interface {
	GetUserByID(id string) (*entities.User, error)
}

type AuthRepositoryStruct struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepositoryStruct {
	return &AuthRepositoryStruct{db}
}

func (db *AuthRepositoryStruct) GetEmailAuthRepository(LoginRequest LoginRequest) (*entities.User, error) {
	// Implement the logic to get a user by ID from the database

	user := &entities.User{}
	if err := db.DB.Preload("Role").First(user, "email = ?", LoginRequest.Email).Error; err != nil {
		return nil, err
	}

	// Check if the user exists
	if user.ID == "" {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (db *AuthRepositoryStruct) GetPhoneCustomerAuthRepository(LoginRequest LoginRequest) (*entities.Customer, error) {
	// Implement the logic to get a user by ID from the database

	user := &entities.Customer{}
	if err := db.DB.First(user, "phone = ?", LoginRequest.Email).Error; err != nil {
		return nil, err
	}

	// Check if the user exists
	if user.ID == "" {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (db *AuthRepositoryStruct) GetTokenAuthRepository(token string) (*entities.User, error) {
	// Implement the logic to get a user by ID from the database

	user := &entities.User{}
	if err := db.DB.First(user, "token = ?", token).Error; err != nil {
		return nil, err
	}

	// Check if the user exists
	if user.ID == "" {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (db *AuthRepositoryStruct) UpdateTokenUserByID(id string, token string) (*entities.User, error) {
	user := &entities.User{}
	if err := db.DB.First(user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	// Check if the user exists
	if user.ID == "" {
		return nil, gorm.ErrRecordNotFound
	}

	// Update the token field
	user.Token = token
	// Save the updated user back to the database
	if err := db.DB.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil

}
