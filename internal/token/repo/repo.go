package tokenRepo

import (
	"time"
	"errors"

	"gorm.io/gorm"

	"selatoz/internal/user/repo"
)

// File handles data logic related to access and refresh tokens

// This defines an AccessToken struct to be used as the access tokens used for authentication
type AccessToken struct {
	gorm.Model
	User      		*userRepo.User  	`json:"user" gorm:"foreignKey:UserID;references:ID"`
	RefreshToken	*RefreshToken		`json:"refresh_token" gorm:"foreignKey:RefreshTokenID;references:ID"`
	UserID			uint					`json:"user_id"gorm:"index"`
	RefreshTokenID	uint					`json:"refresh_token_id",gorm:"index"`
	UserAgent		string				`json:"user_agent",gorm:"index"`
	TokenName		string				`json:"token_name",gorm:"index"`
	TokenString		string				`json:"token_string",gorm:"uniqueIndex"`
	ExpiresAt		time.Time			`json:"expires_at"`
}


// This defines a RefreshToken struct to be used as the refresh tokens used for authentication
type RefreshToken struct {
	gorm.Model
	User      		*userRepo.User  	`json:"user" gorm:"foreignKey:UserID;references:ID"`
	UserID			uint					`json:"user_id,"gorm:"index"`
	UserAgent		string				`json:"user_agent",gorm:"index"`
	TokenName		string				`json:"token_name",gorm:"index"`
	TokenString		string				`json:"token_string",gorm:"uniqueIndex"`
	ExpiresAt		time.Time			`json:"expires_at"`
}

// Repository provides methods for interacting with the profiles in the database
type Repo interface {
	CreateAccessToken(userID uint, refreshTokenID uint, userAgent string, name string, token string, expiresAt time.Time) (*AccessToken, error)
	GetAccessToken(token string) (*AccessToken, error)
	DeleteAccessToken(token string, deleteRefreshToken bool) (error)

	CreateRefreshToken(userID uint, userAgent string, name string, token string, expiresAt time.Time) (*RefreshToken, error)
	GetRefreshToken(token string) (*RefreshToken, error)
	DeleteRefreshToken(token string) (error)
}

// Provides the implementation of the repo
type repo struct {
	db *gorm.DB
}

// NewRepo returns a new instance of the repository with a provided database connection.
func NewRepo(db *gorm.DB) Repo {
	return &repo{db}
}

// CreateAccessToken creates an entry in the access tokens table.
func (r *repo) CreateAccessToken(userID uint, refreshTokenID uint, userAgent string, name string, token string, expiresAt time.Time) (*AccessToken, error) {
	// Create a new personal access token in the database
	at := &AccessToken{
		 UserID:    		userID,
		 RefreshTokenID:	refreshTokenID,
		 UserAgent:			userAgent,
		 TokenName:      	name,
		 TokenString:     token,
		 ExpiresAt: 		expiresAt,
	}

	if err := r.db.Create(at).Error; err != nil {
		 return nil, err
	}

	return at, nil
}

// GetAccessToken returns an access token with the given value
func (r *repo) GetAccessToken(token string) (*AccessToken, error) {
	var at AccessToken
	err := r.db.Where("token_string = ?", token).First(&at).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("token not found")
		}
		return nil, err
	}
	return &at, nil
}

// DeleteAccessToken deletes all matching access token
func (r *repo) DeleteAccessToken(token string, deleteRefreshToken bool) error {
	var at AccessToken

	// Find the matching access token
	err := r.db.Preload("RefreshToken").Where("token_string = ?", token).First(&at).Error
	if err != nil {
		// Handle already deleted
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}

		// Handle other errors
		return err
	}

	// Delete refresh token requested
	// Which deletes the refresh token and all related access tokens
	if (deleteRefreshToken) {
		return r.DeleteRefreshToken(at.RefreshToken.TokenString);
	}

	// Only delete the access token
	return r.db.Where("token_string = ?", token).Delete(&at).Error
}

// CreateRefreshToken creates an entry in the access tokens table.
func (r *repo) CreateRefreshToken(userID uint, userAgent string, name string, token string, expiresAt time.Time) (*RefreshToken, error) {
	// Create a new personal access token in the database
	rt := &RefreshToken{
		 UserID:    	userID,
		 UserAgent:		userAgent,
		 TokenName:    name,
		 TokenString:  token,
		 ExpiresAt: 	expiresAt,
	}

	if err := r.db.Create(rt).Error; err != nil {
		 return nil, err
	}

	return rt, nil
}

// GetRefreshToken returns an access token with the given value
func (r *repo) GetRefreshToken(token string) (*RefreshToken, error) {
	var rt RefreshToken
	err := r.db.Where("token_string = ?", token).First(&rt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("token not found")
		}
		return nil, err
	}
	return &rt, nil
}

// DeleteRefreshToken deletes the matching refresh token along with all referenced access tokens
func (r *repo) DeleteRefreshToken(token string) error {
	var rt RefreshToken

	// Find the matching refresh token
	err := r.db.Where("token_string = ?", token).First(&rt).Error
	if err != nil {
		// Handle already deleted
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}

		// Handle other errors
		return err
	}

	// Delete all access tokens associated with the refresh token
	err = r.db.Where("refresh_token_id = ?", rt.ID).Delete(&AccessToken{}).Error
	if err != nil {
		return err
	}

	// Delete the refresh token
	err = r.db.Delete(&rt).Error
	if err != nil {
		return err
	}

	return nil
}
