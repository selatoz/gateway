package authHttp

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"selatoz/internal/user/svc"
)

type UsersResponse struct {
	Message		string
}

// GetUsersHandler handles the users get request
func GetUsersHandler(userService *userSvc.Svc) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, UsersResponse{Message: "Response from get users"})
	}
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token string `json:"token"`
}

// LoginHandler handles user login requests
func LoginHandler(userService userSvc.Svc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bind the request body to a LoginRequest struct
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Authenticate the user
		_, err := userService.Authenticate(req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		// // Generate a JWT token for the user
		// token, err = userService.GenerateAccessToken(user.ID)
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		// 	return
		// }

		// Return the token in a LoginResponse struct
		c.JSON(http.StatusOK, LoginResponse{Token: "token"})
	}
}