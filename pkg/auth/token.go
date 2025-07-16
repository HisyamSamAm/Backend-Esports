package auth

import (
	"embeck/model"
	"errors"
	"os"
	"time"

	"aidanwoods.dev/go-paseto"
)

var (
	secretKey = []byte(getSecretKey())
)

// getSecretKey returns the secret key from environment or default
func getSecretKey() string {
	key := os.Getenv("PASETO_SECRET_KEY")
	if key == "" {
		// Default key for development - CHANGE THIS IN PRODUCTION!
		key = "embeck-secret-key-change-this-in-production-32-chars"
	}
	// Ensure key is exactly 32 bytes for PASETO v2
	if len(key) < 32 {
		key = key + "000000000000000000000000000000000000"
	}
	return key[:32]
}

// GenerateToken creates a new PASETO token for a user
func GenerateToken(user *model.User) (string, error) {
	// Create a new PASETO token
	token := paseto.NewToken()

	// Set token to be a local token (encrypted)
	token.SetIssuedAt(time.Now())
	token.SetExpiration(time.Now().Add(24 * time.Hour)) // Token expires in 24 hours
	token.SetNotBefore(time.Now())

	// Set custom claims
	claims := model.TokenClaims{
		UserID:   user.ID.Hex(),
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		IssuedAt: time.Now().Unix(),
		ExpireAt: time.Now().Add(24 * time.Hour).Unix(),
	}

	token.SetString("user_id", claims.UserID)
	token.SetString("username", claims.Username)
	token.SetString("email", claims.Email)
	token.SetString("role", claims.Role)

	// Create symmetric key from secret
	symmetricKey, err := paseto.V4SymmetricKeyFromBytes(secretKey)
	if err != nil {
		return "", err
	}

	// Encrypt the token
	encrypted := token.V4Encrypt(symmetricKey, nil)

	return encrypted, nil
}

// ValidateToken validates and parses a PASETO token
func ValidateToken(tokenString string) (*model.TokenClaims, error) {
	// Create symmetric key from secret
	symmetricKey, err := paseto.V4SymmetricKeyFromBytes(secretKey)
	if err != nil {
		return nil, err
	}

	// Parse the token
	parser := paseto.NewParser()
	token, err := parser.ParseV4Local(symmetricKey, tokenString, nil)
	if err != nil {
		return nil, errors.New("invalid token format")
	}

	// Check if token is expired
	expiration, err := token.GetExpiration()
	if err != nil {
		return nil, errors.New("invalid expiration in token")
	}
	if time.Now().After(expiration) {
		return nil, errors.New("token has expired")
	}

	// Check if token is not yet valid
	notBefore, err := token.GetNotBefore()
	if err != nil {
		return nil, errors.New("invalid not_before in token")
	}
	if time.Now().Before(notBefore) {
		return nil, errors.New("token not yet valid")
	}

	// Extract claims
	claims := &model.TokenClaims{}

	userID, err := token.GetString("user_id")
	if err != nil {
		return nil, errors.New("invalid user_id in token")
	}
	claims.UserID = userID

	username, err := token.GetString("username")
	if err != nil {
		return nil, errors.New("invalid username in token")
	}
	claims.Username = username

	email, err := token.GetString("email")
	if err != nil {
		return nil, errors.New("invalid email in token")
	}
	claims.Email = email

	role, err := token.GetString("role")
	if err != nil {
		return nil, errors.New("invalid role in token")
	}
	claims.Role = role

	issuedAt, err := token.GetIssuedAt()
	if err == nil {
		claims.IssuedAt = issuedAt.Unix()
	}

	claims.ExpireAt = expiration.Unix()

	return claims, nil
}
