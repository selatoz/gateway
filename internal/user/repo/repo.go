package userRepo

import (
	"gorm.io/gorm"
)

// File handles business logic related to the user
// Table name is the plural of the name of the User struct (users)

// This defines a User struct that represents a user record
type User struct {
	gorm.Model
	Email string 		`json:"email",gorm:"uniqueIndex;index"`
	Password string	`json:"password"`
}

// Repository provides methods for interacting with the profiles in the database
type Repo interface {
	GetByEmail(email string) (*User, error)
	GetById(userID uint) (*User, error)
}

type repo struct {
	db *gorm.DB
}

// NewRepo returns a new instance of the repository with a provided database connection.
func NewRepo(db *gorm.DB) Repo {
	return &repo{db}
}

// GetByEmail returns a user record based on the provided email
func (r *repo) GetByEmail(email string) (*User, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetById returns a user record based on the provided id
func (r *repo) GetById(userID uint) (*User, error) {
	var user User
	err := r.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
