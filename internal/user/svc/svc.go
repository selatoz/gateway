package userSvc

import (
	"errors"

	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"

	"selatoz/internal/user/repo"
)

// Define errors
var (
	ErrInvalidPassword 			= errors.New("invalid password")
	ErrFailedToHashPassword 	= errors.New("failed to hash password")
	ErrFailedToCreateUser 		= errors.New("failed to create user")
	ErrFailedToGenerateToken 	= errors.New("failed to generate token")
)

// Svc is an interface for defining the methods that the user service will provide.
type Svc interface {
	Login(email string, password string) (*userRepo.User, error)
	Register(email string, password string) (*userRepo.User, error)
	GetById(id uint) (*userRepo.User, error)
	GetByEmail(email string) (*userRepo.User, error)
	// Add more methods here as needed
}

// svc is an implementation of the UserSvc interface that handles the business logic for user-related operations.
type svc struct {
	repo userRepo.Repo
}

// NewSvc creates a new instance of svc and returns it as a Svc interface.
func NewSvc(db *gorm.DB) Svc {
	repo := userRepo.NewRepo(db)
	return &svc{repo}
}

// GetById returns a user record based on the given id
func (s *svc) GetById(id uint) (*userRepo.User, error) {
	return s.repo.GetById(id)
}

// GetByEmail returns a user record based on the given email
func (s *svc) GetByEmail(email string) (*userRepo.User, error) {
	return s.repo.GetByEmail(email)
}

// Login checks if a user exists with the given credentials
func (s *svc) Login(email string, password string) (*userRepo.User, error) {
	// Get the user by email
	u, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	// Check password matches
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return nil, ErrInvalidPassword
	}

	// Implement logic for creating a new user
	return u, nil
}

// Register creates a new user with the given information
func (s *svc) Register(email string, password string) (*userRepo.User, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrFailedToHashPassword
	}

	// Create a new user
	u, err := s.repo.NewUser(email, string(hashedPassword))
	if err != nil {
		return nil, ErrFailedToCreateUser
	}

	// Handle return
	return u, nil
}

// Check if the provided password matches the password hash
func passwordsMatch(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}