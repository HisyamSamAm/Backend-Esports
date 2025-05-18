package model

type User struct {
	ID       string `json:"id" bson:"id,omitempty"`
	Username string `json:"username" bson:"username" validate:"required"`
	Role     string `json:"role" bson:"role"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}
type UserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
