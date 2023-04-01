package routes

import (
	"fmt"
	"net/http"
	"hash/fnv"
	"encoding/base64"
	
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/selatoz/gateway/api/auth"
	"github.com/selatoz/gateway/middleware/auth"
	"github.com/selatoz/gateway/internal/user/svc"
	"github.com/selatoz/gateway/internal/token/svc"
)

// Context represents the extended gin.Context type
type Context struct {
	*gin.Context
	UserService		userSvc.Svc
}

// Route represents a single API route.
type Route struct {
	Method  		string
	Path    		string
	Handler 		gin.HandlerFunc
	Middleware	[]gin.HandlerFunc
}

// RouteGroup represents a group of routes, 
// grouped based on the middleware they use.
type RouteGroup struct {
	Middleware []gin.HandlerFunc
	Routes     []Route
}

// Routes represents a collection of API routes.
type Routes []Route

// Initializes the router object with the routes
func NewRoutes(router *gin.Engine, db *gorm.DB) {
	userService := userSvc.NewSvc(db)
	tokenService := tokenSvc.NewSvc(db)

	// Define API routes
	apiRoutes := Routes{
		{
			Method:  	"GET",
			Path:    	"/users",
			Handler: 	authHttp.GetUsersHandler(userService),
			Middleware: []gin.HandlerFunc{mwauth.Authorize(tokenService)},
		},
		{
			Method:  "POST",
			Path:    "/user/logout",
			Handler: authHttp.LogoutHandler(tokenService),
			Middleware: []gin.HandlerFunc{mwauth.Authorize(tokenService)},
		},
		{
			Method:  "POST",
			Path:    "/auth/login",
			Handler: authHttp.LoginHandler(userService, tokenService),
			Middleware: nil,
		},
		{
			Method:  "POST",
			Path:    "/auth/register",
			Handler: authHttp.RegisterHandler(userService, tokenService),
			Middleware: nil,
		},
		{
			Method:  "POST",
			Path:    "/auth/refresh-access",
			Handler: authHttp.RefreshAccessHandler(tokenService),
			Middleware: nil,
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
			Middleware: nil,
		},
	}

	// Register API routes
	apiRoutes.RegisterRoute(router)

	// Register TEST routes
	testRoutes.RegisterRoute(router)
}

func getMiddlewareHash(mw []gin.HandlerFunc) string {
	// Handle nil
	if mw != nil {
		h := fnv.New32a()
		for _, f := range mw {
			 h.Write([]byte(fmt.Sprintf("%p", f)))
		}
		return base64.URLEncoding.EncodeToString(h.Sum(nil))
	}

	return "no_mw"
}

// Register registers the API routes with the provided Gin router.
func (routes Routes) RegisterRoute(router *gin.Engine) {
	// Group routes by middleware
	groups := make(map[string]RouteGroup)
	for _, route := range routes {
		key := getMiddlewareHash(route.Middleware)
		group, ok := groups[key]
		if !ok {
			 group = RouteGroup{
				  Middleware: route.Middleware,
			 }
		}
		group.Routes = append(group.Routes, route)
		groups[key] = group
  }
  
  // Register routes with their respective middleware groups
	for _, group := range groups {
		chain := router.Group("/")
		chain.Use(group.Middleware...)
		for _, route := range group.Routes {
			chain.Handle(route.Method, route.Path, route.Handler)
		}
	}
}





