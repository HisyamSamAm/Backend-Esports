package auth

import (
	"embeck/model"
	"encoding/hex"
	"errors"
	"os"
	"time"

	"aidanwoods.dev/go-paseto"
)

// GenerateToken creates a new PASETO token for a user using public-key signing.
func GenerateToken(user *model.User) (string, error) {
	// Get the hex-encoded private key from the environment variables.
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	if privateKeyHex == "" {
		return "", errors.New("PRIVATE_KEY environment variable not set")
	}

	// Decode the hex string into bytes.
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", errors.New("failed to decode private key")
	}

	// Create a Paseto V4 asymmetric private key from the decoded bytes.
	privateKey, err := paseto.NewV4AsymmetricSecretKeyFromBytes(privateKeyBytes)
	if err != nil {
		return "", errors.New("failed to create paseto private key")
	}

	// Create a new token.
	token := paseto.NewToken()
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(24 * time.Hour)) // Token is valid for 24 hours.

	// Add custom claims to the token.
	token.SetString("user_id", user.ID.Hex())
	token.SetString("username", user.Username)
	token.SetString("email", user.Email)
	token.SetString("role", user.Role)

	// Sign the token with the private key to create a public token.
	signedToken := token.V4Sign(privateKey, nil)

	return signedToken, nil
}

// ValidateToken validates and parses a PASETO token using the public key.
func ValidateToken(tokenString string) (*model.TokenClaims, error) {
	// Get the hex-encoded public key from the environment variables.
	publicKeyHex := os.Getenv("PUBLIC_KEY")
	if publicKeyHex == "" {
		return nil, errors.New("PUBLIC_KEY environment variable not set")
	}

	// Decode the hex string into bytes.
	publicKeyBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return nil, errors.New("failed to decode public key")
	}

	// Create a Paseto V4 asymmetric public key from the decoded bytes.
	publicKey, err := paseto.NewV4AsymmetricPublicKeyFromBytes(publicKeyBytes)
	if err != nil {
		return nil, errors.New("failed to create paseto public key")
	}

	// Create a new Paseto parser to verify the token.
	parser := paseto.NewParser()
	// Add rules to validate standard claims (e.g., token has not expired).
	parser.AddRule(paseto.NotExpired())

	// Parse and validate the token using the public key.
	token, err := parser.ParseV4Public(publicKey, tokenString, nil)
	if err != nil {
		return nil, errors.New("invalid token or signature")
	}

	// Extract custom claims from the token.
	claims := &model.TokenClaims{}
	if err := token.Get("user_id", &claims.UserID); err != nil {
		return nil, err
	}
	if err := token.Get("username", &claims.Username); err != nil {
		return nil, err
	}
	if err := token.Get("email", &claims.Email); err != nil {
		return nil, err
	}
	if err := token.Get("role", &claims.Role); err != nil {
		return nil, err
	}

	// Extract standard time-based claims.
	issuedAt, err := token.GetIssuedAt()
	if err == nil {
		claims.IssuedAt = issuedAt.Unix()
	}
	expiration, err := token.GetExpiration()
	if err == nil {
		claims.ExpireAt = expiration.Unix()
	}

	return claims, nil
}
