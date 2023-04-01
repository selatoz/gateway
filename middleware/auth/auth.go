package mwauth

import (
	"fmt"
	"strings"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/selatoz/gateway/pkg/cfglib"
	"github.com/selatoz/gateway/validation/http"
	"github.com/selatoz/gateway/internal/token/svc"
)

// Define constants
const (
	// Errors
	ErrNoAuthorization = "Missing authorization"

	// Headers
	HeaderAuthorization 				= "Authorization"
	HeaderRefreshAuthorization 	= "Refresh-Authorization"
	HeaderChallengeAuthorization	= "WWW-Authenticate"

	// Header values
	ChallengeExpiredAccessToken 	= "Bearer realm=\"%s\",error=\"access_token_expired\""
)

type AuthContext struct {
	UserID			uint
	AccessToken		string
}

// Middleware is a function that wraps an gin.HandlerFunc and provides some extra functionality.
type NewAuthMiddleware func(gin.HandlerFunc) gin.HandlerFunc

// authMiddleware is a middleware that requires an authorization header with a valid token to access a route.
func Authorize(tokenService tokenSvc.Svc) gin.HandlerFunc {
	return func(c *gin.Context) {
		at := c.GetHeader(HeaderAuthorization)
		if at == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, validHttp.ErrorResponse{Error: ErrNoAuthorization})
			return
		}

		// Validate token
		isLoggingOut := (c.Request.Method == http.MethodPost && c.Request.URL.Path == "/user/logout")
		uid, _, err := tokenService.ValidateToken(at, isLoggingOut)
		if err != nil {
			// Check if the error is due to an expired token
			if err.Error() == tokenSvc.ErrTokenExpired {
				// Challenge the client to send the refresh token
				c.Writer.Header().Set(HeaderChallengeAuthorization, fmt.Sprintf(ChallengeExpiredAccessToken, cfglib.DefaultConf.AppName))
				c.AbortWithStatusJSON(http.StatusUnauthorized, validHttp.ErrorResponse{Error: err.Error()})
				return
		  }

			c.AbortWithStatusJSON(http.StatusUnauthorized, validHttp.ErrorResponse{Error: err.Error()})
			return
		}

		// Create auth context
		authContext := AuthContext{
			UserID:       uid,
			AccessToken:  strings.TrimPrefix(at, "Bearer "),
		}

		// Add auth context to request context
		c.Set("auth", &authContext)
		c.Next()
	}
}