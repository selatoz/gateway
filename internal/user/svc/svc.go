package userSvc

import (
	"gorm.io/gorm"

	"selatoz/internal/user/repo"
)

// Svc is an interface for defining the methods that the user service will provide.
type Svc interface {
	Authenticate(email string, password string) (*userRepo.User, error)
	GetUserByID(id int) (*userRepo.User, error)
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

// CreateUser creates a new user with the given information.
func (s *svc) Authenticate(email string, password string) (*userRepo.User, error) {
	var u = userRepo.User{
		Email: 	"test@email.com",
	}

	// Implement logic for creating a new user
	return &u, nil
}

// GetUserByID retrieves a user with the given ID.
func (s *svc) GetUserByID(id int) (*userRepo.User, error) {
	// Implement logic for retrieving a user by ID
	return nil, nil
}