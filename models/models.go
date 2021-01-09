package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwt"
)

// Credentials represents the user credentials
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Account represents the account model
type Account struct {
	ID         string    `json:"id" db:"id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	ModifiedAt time.Time `json:"modified_at" db:"modified_at"`
}

// JWT represents a jwt wrapper
type JWT jwt.Token

// NewJWT creates a new JWT wrapper type
func NewJWT() JWT {
	token := JWT(jwt.New())
	_ = token.Set(jwt.IssuedAtKey, time.Now().UTC().Unix())
	_ = token.Set(jwt.IssuerKey, "account-api-srv")
	_ = token.Set(jwt.JwtIDKey, uuid.New())
	_ = token.Set(jwt.ExpirationKey, time.Now().Add(time.Minute*2).UTC().Unix())
	return token
}
