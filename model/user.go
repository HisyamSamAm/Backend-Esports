package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user entity
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty" example:"64f123abc456def789012345" description:"ID unik user"`
	Username  string             `bson:"username" json:"username" example:"userbaru123" description:"Nama pengguna untuk login"`
	Email     string             `bson:"email" json:"email" example:"user.example@example.com" description:"Alamat email user"`
	Password  string             `bson:"password" json:"-" description:"Password yang telah di-hash (tidak ditampilkan di response)"`
	Role      string             `bson:"role" json:"role" example:"user" description:"Peran user: admin atau user"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at" example:"2025-07-16T07:28:37.016Z" description:"Waktu pembuatan user"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at" example:"2025-07-16T07:28:37.016Z" description:"Waktu terakhir diupdate"`
}

// RegisterRequest represents request body for user registration
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50" example:"userbaru123" description:"Nama pengguna (3-50 karakter)"`
	Email    string `json:"email" validate:"required,email" example:"user.baru@example.com" description:"Alamat email yang valid"`
	Password string `json:"password" validate:"required,min=6" example:"passwordAman123" description:"Password minimal 6 karakter"`
}

// LoginRequest represents request body for user login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user.example@example.com" description:"Alamat email untuk login"`
	Password string `json:"password" validate:"required" example:"passwordAman123" description:"Password user"`
}

// UpdateUserRequest represents request body for updating user data
type UpdateUserRequest struct {
	Username string `json:"username,omitempty" validate:"omitempty,min=3,max=50" example:"usernameUpdate" description:"Nama pengguna baru (opsional)"`
	Email    string `json:"email,omitempty" validate:"omitempty,email" example:"email.update@example.com" description:"Email baru (opsional)"`
	Role     string `json:"role,omitempty" validate:"omitempty,oneof=user admin" example:"admin" description:"Role baru: user atau admin (opsional)"`
}

// AuthResponse represents response for authentication operations
type AuthResponse struct {
	Message string `json:"message" example:"Login successful" description:"Pesan konfirmasi"`
	Token   string `json:"token,omitempty" example:"v2.local.xxx" description:"PASETO token untuk autentikasi"`
	Role    string `json:"role,omitempty" example:"user" description:"Peran user"`
	UserID  string `json:"user_id,omitempty" example:"64f123abc456def789012345" description:"ID user yang login"`
}

// UserResponse represents response for user operations (without sensitive data)
type UserResponse struct {
	Message string `json:"message" example:"User registered successfully" description:"Pesan konfirmasi"`
	UserID  string `json:"user_id,omitempty" example:"64f123abc456def789012345" description:"ID user yang dibuat"`
}

// UsersListResponse represents response for get all users
type UsersListResponse struct {
	Message string        `json:"message" example:"Users retrieved successfully" description:"Pesan konfirmasi"`
	Data    []UserProfile `json:"data" description:"Daftar semua user"`
	Total   int           `json:"total" example:"5" description:"Total jumlah user"`
}

// UserProfile represents user profile data (without password)
type UserProfile struct {
	ID        primitive.ObjectID `json:"_id,omitempty" example:"64f123abc456def789012345"`
	Username  string             `json:"username" example:"userbaru123"`
	Email     string             `json:"email" example:"user.example@example.com"`
	Role      string             `json:"role" example:"user"`
	CreatedAt time.Time          `json:"created_at" example:"2025-07-16T07:28:37.016Z"`
	UpdatedAt time.Time          `json:"updated_at" example:"2025-07-16T07:28:37.016Z"`
}

// TokenClaims represents the claims stored in PASETO token
type TokenClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	IssuedAt int64  `json:"iat"`
	ExpireAt int64  `json:"exp"`
}
