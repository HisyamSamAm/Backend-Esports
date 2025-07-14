package model

type User struct {
	ID       string `json:"id" bson:"id,omitempty"`
	Username string `json:"username" bson:"username" validate:"required"`
	Email    string `json:"email" bson:"email"`
	Role     string `json:"role" bson:"role"`
	Password string `json:"password" bson:"password"`
}
type UserLogin struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Role     string `json:"role" bson:"role"`
}

type RegisterRequest struct {
	Username string `json:"username" bson:"username" validate:"required"`
	Email    string `json:"email" bson:"email" validate:"required,email"`
	Password string `json:"password" bson:"password" validate:"required"`
}

type Payload struct {
	User string `json:"user"`
	Role string `json:"role"`
}
