package model

type User struct {
	ID       string `json:"id" bson:"id,omitempty" example:"user-001"`
	Username string `json:"username" bson:"username" validate:"required" example:"john_doe"`
	Role     string `json:"role" bson:"role" example:"user"`
	Email    string `json:"email" bson:"email" example:"john@example.com"`
	Password string `json:"password" bson:"password" example:"password123"`
}

type UserInput struct {
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"password123"`
}

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Username string `json:"username" example:"john_doe" validate:"required"`
	Email    string `json:"email" example:"john@example.com" validate:"required"`
	Password string `json:"password" example:"password123" validate:"required"`
	Role     string `json:"role" example:"user"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username" example:"john_doe" validate:"required"`
	Password string `json:"password" example:"password123" validate:"required"`
}

// UserResponse represents the response for user operations
type UserResponse struct {
	Status  int    `json:"status" example:"200"`
	Message string `json:"message" example:"Berhasil nambahin data user bre!"`
	Data    User   `json:"data"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Status  int    `json:"status" example:"200"`
	Message string `json:"message" example:"Berhasil login sebagai user"`
	Data    User   `json:"data"`
}
