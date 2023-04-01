package validHttp

// Generic types
type ErrorResponse struct {
	Error		string	`json:"error"`
}

type SuccessResponse struct {
	Message	string	`json:"message"`
}

type EmptyRequest struct {}

// Types related to auth
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required,min=6"`
}

type RefreshAccessRequest struct {}

// Types related to profile
type UpdateUserProfileRequest struct {
	Message string `json:"message"`
}