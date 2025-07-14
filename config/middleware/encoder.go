package middleware

import (
	"EMBECK/model"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"aidanwoods.dev/go-paseto"
)

func EncodeWithRoleHours(role, username string, hours int64) (string, error) {
	privatekey := os.Getenv("PRIVATEKEY")
	token := paseto.NewToken()
	// Set metadata: waktu pembuatan, masa berlaku, dll
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(time.Duration(hours) * time.Hour))
	token.SetString("user", username)
	token.SetString("role", role)
	key, err := paseto.NewV4AsymmetricSecretKeyFromHex(privatekey)
	return token.V4Sign(key, nil), err
}

func Decoder(tokenstr string) (payload model.Payload, err error) {
	publickey := os.Getenv("PUBLIC_KEY")

	pubKey, err := paseto.NewV4AsymmetricPublicKeyFromHex(publickey)
	if err != nil {
		fmt.Println("Decode NewV4AsymmetricPublicKeyFromHex : ", err)
	}

	parser := paseto.NewParser()
	token, err := parser.ParseV4Public(pubKey, tokenstr, nil)
	if err != nil {
		fmt.Println("Decode ParseV4Public : ", err)
	} else {
		json.Unmarshal(token.ClaimsJSON(), &payload)
	}

	return payload, err
}