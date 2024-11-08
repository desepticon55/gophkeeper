package model

import (
	"github.com/golang-jwt/jwt/v4"
)

const (
	CredentialsSecretType = "CREDENTIALS"
	TextSecretType        = "TEXT"
	CardSecretType        = "CARD"
	BinarySecretType      = "BINARY"
)

// JWT claims
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// User domain model
type User struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

// Secret domain model
type Secret struct {
	Name     string `json:"name"`
	Content  []byte `json:"content"`
	Username string `json:"-"`
	Type     string `json:"type"`
	Version  int64  `json:"-"`
}

// Card requisites
type Card struct {
	Number string `json:"number"`
	Date   string `json:"date"`
	Code   string `json:"code"`
	Holder string `json:"holder"`
}
