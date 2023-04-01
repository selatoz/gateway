package authHttp

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/selatoz/gateway/validation/http"
	"github.com/selatoz/gateway/middleware/auth"
	"github.com/selatoz/gateway/internal/user/svc"
	"github.com/selatoz/gateway/internal/token/svc"
)

// Set constants
const (
	// Errors
	ErrInvalidCredentials 				= "Invalid credentials"
	ErrMissingContext						= "Missing context"
	ErrFailedToGenerateRefreshToken 	= "Failed to generate <r> token"
	ErrFailedToGenerateAccessToken	= "Failed to generate <a> token"
	ErrUserAlreadyExists					= "User already exists"
)

type UsersResponse struct {
	Message		string
}

// GetUsersHandler handles the users get request
func GetUsersHandler(userService userSvc.Svc) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, UsersResponse{Message: "Response from get users"})
	}
}

// LoginHandler handles user login request
func LoginHandler(userService userSvc.Svc, tokenService tokenSvc.Svc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bind the request body to a LoginRequest struct
		var req validHttp.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get the user agent
		userAgent := c.Request.UserAgent()

		// Authenticate the user
		u, err := userService.Login(req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, validHttp.ErrorResponse{Error: ErrInvalidCredentials})
			return
		}

		// Generate the tokens
		at, rt, err := tokenService.GenerateTokens(u.ID, userAgent)
		if err != nil {
			c.JSON(http.StatusInternalServerError, validHttp.ErrorResponse{Error: err.Error()})
			return
		}

		// Set the token headers
		c.Header(mwauth.HeaderAuthorization, "Bearer "+at.TokenString)
		c.Header(mwauth.HeaderRefreshAuthorization, rt.TokenString)

		// Return the token in a LoginResponse struct
		c.JSON(http.StatusOK, validHttp.SuccessResponse{Message: "Login success"})
	}
}

// LogoutHandler handles user logout request
func LogoutHandler(tokenService tokenSvc.Svc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read auth context
		authCtx := c.MustGet("auth").(*mwauth.AuthContext)
		if authCtx == nil {
			c.JSON(http.StatusBadRequest, validHttp.ErrorResponse{Error: ErrMissingContext})
			return
		}

		// Bind the request body
		var req validHttp.EmptyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusInternalServerError, validHttp.ErrorResponse{Error: err.Error()})
			return
		}

		err := tokenService.DeleteAccessToken(authCtx.AccessToken, true)
		if err != nil {
			c.JSON(http.StatusInternalServerError, validHttp.ErrorResponse{Error: err.Error()})
			return
		}

		// Handle success
		c.JSON(http.StatusOK, validHttp.SuccessResponse{Message: "Login success"})
	}
}

// RegisterHandler handles user registration request
func RegisterHandler(userService userSvc.Svc, tokenService tokenSvc.Svc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bind the request body to a RegisterRequest struct
		var req validHttp.RegisterRequest
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get the user agent
		userAgent := c.Request.UserAgent()

		// Check if user with the same email exists
		if _, err := userService.GetByEmail(req.Email); err == nil {
			c.JSON(http.StatusConflict, validHttp.ErrorResponse{Error: err.Error()})
			return
		}

		// Register the user
		u, err := userService.Register(req.Email, req.Password)

		// Generate the tokens
		at, rt, err := tokenService.GenerateTokens(u.ID, userAgent)
		if err != nil {
			c.JSON(http.StatusInternalServerError, validHttp.ErrorResponse{Error: err.Error()})
			return
		}

		// Set the token headers
		c.Header(mwauth.HeaderAuthorization, "Bearer "+at.TokenString)
		c.Header(mwauth.HeaderRefreshAuthorization, rt.TokenString)

		// Return the token in a RegisterResponse struct
		c.JSON(http.StatusOK, validHttp.SuccessResponse{Message: "Register success"})
	}
}

// RefreshTokenHandler handles access token refreshing
func RefreshAccessHandler(tokenService tokenSvc.Svc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bind the request body to a RegisterRequest struct
		var req validHttp.RefreshAccessRequest
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, validHttp.ErrorResponse{Error: err.Error()})
			return
		}

		// Get the user agent
		userAgent := c.Request.UserAgent()

		// Get the token from the header
		currRt := c.GetHeader("Authorization")
		currRt = strings.TrimPrefix(currRt, "Bearer ")

		// Validate the refresh token
		userID, _, err := tokenService.ValidateToken(currRt, false)
		if err != nil {
			 c.JSON(http.StatusUnauthorized, validHttp.ErrorResponse{Error: err.Error()})
			 return
		}

		// Delete the old tokens
		if err := tokenService.DeleteRefreshToken(currRt); err != nil {
			 c.JSON(http.StatusInternalServerError, validHttp.ErrorResponse{Error: err.Error()})
			 return
		}

		// Issue new tokens
		at, rt, err := tokenService.GenerateTokens(userID, userAgent)
		if err != nil {
			c.JSON(http.StatusInternalServerError, validHttp.ErrorResponse{Error: err.Error()})
			return
		}

		// Add the new tokens to the response header
		c.Header(mwauth.HeaderAuthorization, "Bearer " + at.TokenString)
		c.Header(mwauth.HeaderRefreshAuthorization, rt.TokenString)

		// Handle success
		c.JSON(http.StatusOK, validHttp.SuccessResponse{Message: "Refresh success"})
	}
}