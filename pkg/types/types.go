package types

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	UUID  int    `json:"uuid"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

type ContextKeyType string
