package tokenSvc

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
	"github.com/dgrijalva/jwt-go"

	"selatoz/pkg/cfglib"
	"selatoz/internal/token/repo"
	"selatoz/internal/user/repo"
)

// Define constants
const (
	// Errors
	ErrTokenExpired 			= "token expired"
	ErrTokenInvalid			= "token invalid"
	ErrTokenUserInvalid		= "token user invalid"
	ErrTokenClaimsInvalid 	= "token claims invalid"

	// These names are used to determine if the token is an access token or a refresh token
	AccessTokenName	= "jwt_access"
	RefreshTokenName	= "jwt_refresh"
)

// Svc is an interface for defining the methods that the user service will provide.
type Svc interface {
	GetAccessToken(token string) (*tokenRepo.AccessToken, error)
	GenerateTokens(userID uint, userAgent string) (*tokenRepo.AccessToken, *tokenRepo.RefreshToken, error)
	GenerateAccessToken(userID uint, refreshTokenID uint, userAgent string) (*tokenRepo.AccessToken, error)
	GenerateRefreshToken(userID uint, userAgent string) (*tokenRepo.RefreshToken, error)
	DeleteRefreshToken(token string) (error)
	DeleteAccessToken(token string, deleteRelatedRefreshToken bool) (error)
	ValidateToken(token string, allowExpired bool) (uint, string, error)
	// Add more methods here as needed
}

// svc is an implementation of the UserSvc interface that handles the business logic for user-related operations.
type svc struct {
	repo 			tokenRepo.Repo
	userRepo		userRepo.Repo
}

// NewSvc creates a new instance of svc and returns it as a Svc interface.
func NewSvc(db *gorm.DB) Svc {
	ur := userRepo.NewRepo(db)
	tr := tokenRepo.NewRepo(db)

	return &svc{
		repo: tr,
		userRepo: ur,
	}
}

// GetAccessToken provides access to get the access token record
func (s *svc) GetAccessToken(token string) (*tokenRepo.AccessToken, error) {
	return s.repo.GetAccessToken(token);
}

/*
 * This method generates both the access and the refresh tokens,
 * stores them in the appropriate database tables,
 * and returns them.
 * 
 * @userID - the id of the user to which the tokens will belong
*/
func (s *svc) GenerateTokens(userID uint, userAgent string) (*tokenRepo.AccessToken, *tokenRepo.RefreshToken, error) {
	// Generate the refresh token first, as it is needed to make the access token
	rt, err := s.GenerateRefreshToken(userID, userAgent)
	if err != nil {
		return nil, nil, err
	}

	// Generate the access token using the refresh token
	at, err := s.GenerateAccessToken(userID, rt.ID, userAgent)
	if err != nil {
		return nil, nil, err
	}

	// Handle success
	return at, rt, nil
}

/* 
 * This method generates a JSON Web Token (JWT) with a payload that includes the user ID and a short expiration time.
 * The token is signed using a secret key provided in the application configuration.
 */
func (s *svc) GenerateAccessToken(userID uint, refreshTokenID uint, userAgent string) (*tokenRepo.AccessToken, error) {
	// Load configs
	secretKey := cfglib.DefaultConf.AppSecret
	tokenExpiration := time.Duration(cfglib.DefaultConf.TokenExpAccess)
	expiresAt := time.Unix(time.Now().Add(tokenExpiration * time.Hour).Unix(), 0)

	// Generate a new JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	
	// Set claims for the token
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user_id"] = userID
	claims["name"] = AccessTokenName
	claims["exp"] = expiresAt

	// Sign the token
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	// Store in the database
	return s.repo.CreateAccessToken(userID, refreshTokenID, userAgent, AccessTokenName, tokenString, expiresAt)
}

/* 
 * This method generates a JSON Web Token (JWT) with a payload that includes the user ID and a long expiration time.
 * The token is signed using a secret key provided in the application configuration.
 */
 func (s *svc) GenerateRefreshToken(userID uint, userAgent string) (*tokenRepo.RefreshToken, error) {
	// Load configs
	secretKey := cfglib.DefaultConf.AppSecret
	tokenExpiration := time.Duration(cfglib.DefaultConf.TokenExpRefresh)
	expiresAt := time.Unix(time.Now().Add(tokenExpiration * time.Hour).Unix(), 0)

	// Generate a new JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	
	// Set claims for the token
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user_id"] = userID
	claims["name"] = RefreshTokenName
	claims["exp"] = expiresAt

	// Sign the token
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	// Store in the database
	return s.repo.CreateRefreshToken(userID, userAgent, RefreshTokenName, tokenString, expiresAt)
}

/*
 * This method removes the access token, and the related refresh token if specified
*/
func (s *svc) DeleteAccessToken(token string, deleteRelatedRefreshToken bool) (error) {
	return s.repo.DeleteAccessToken(token, deleteRelatedRefreshToken)
}

/*
 * This method removes the given refresh token, along with all related access tokens
*/
func (s *svc) DeleteRefreshToken(token string) (error) {
	return s.repo.DeleteRefreshToken(token)
}

/*
 * This method validates a JWT and returns the user ID if the token is valid.
 * Returns <userId, tokenString, error>
 */
func (s *svc) ValidateToken(token string, allowExpired bool) (uint, string, error) {
	// Load configs
	secretKey := cfglib.DefaultConf.AppSecret

	// Remove "Bearer " prefix from token
	token = strings.TrimPrefix(token, "Bearer ")

	// Validate the token string
	jt, err := jwt.Parse(token, func(jt *jwt.Token) (interface{}, error) {
		if _, ok := jt.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(ErrTokenInvalid)
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, "", err
	}

	if claims, ok := jt.Claims.(jwt.MapClaims); ok && jt.Valid {
		// Handle invalid claims
		uid, ok := claims["user_id"].(float64)
		if !ok {
			return 0, "", errors.New(ErrTokenClaimsInvalid + " <user>")
		}
		exp, ok := claims["exp"].(string)
		if !ok {
			return 0, "", errors.New(ErrTokenClaimsInvalid + " <exp>")
		}
		name, ok := claims["name"].(string)
		if !ok {
			return 0, "", errors.New(ErrTokenClaimsInvalid + " <name>")
		}

		// Handle expired token
		if !allowExpired {
			expTime, err := time.Parse(time.RFC3339, exp)
			if err != nil {
				return 0, "", err
			}

			if time.Now().UTC().After(expTime) {
				return 0, "", errors.New(ErrTokenExpired)
			}
		}

		// Handle token mismatch
		user, err := s.userRepo.GetById(uint(uid))
		if (err != nil || user.ID != uint(uid)) {
			return 0, "", errors.New(ErrTokenUserInvalid)
		}

		return user.ID, name, nil
	}

	return 0, "", errors.New(ErrTokenInvalid)
}