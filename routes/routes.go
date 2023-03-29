package routes

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"selatoz/api/auth"
	"selatoz/internal/user/svc"
)

// Context represents the extended gin.Context type
type Context struct {
	*gin.Context
	UserService		userSvc.Svc
}

// Route represents a single API route.
type Route struct {
	Method  string
	Path    string
	Handler gin.HandlerFunc
}

// Routes represents a collection of API routes.
type Routes []Route

// Initializes the router object with the routes
func Init(router *gin.Engine, db *gorm.DB) {
	userService := userSvc.NewSvc(db)

	// Define API routes
	apiRoutes := Routes{
		{
			Method:  "GET",
			Path:    "/users",
			Handler: authHttp.GetUsersHandler(&userService),
		},
		{
			Method:  "POST",
			Path:    "/users",
			Handler: func(c *gin.Context) {
				c.String(http.StatusOK, "POST users")
			},
		},
	}

	// Define TEST routes
	testRoutes := Routes{
		{
			Method:  "GET",
			Path:    "/ping",
			Handler: func(c *gin.Context) {
				c.String(http.StatusOK, "pong")
			},
		},
	}

	// Register API routes
	apiRoutes.Register(router)

	// Register TEST routes
	testRoutes.Register(router)
}

// Register registers the API routes with the provided Gin router.
func (routes Routes) Register(router *gin.Engine) {
	for _, route := range routes {
		router.Handle(route.Method, route.Path, route.Handler)
	}
}





