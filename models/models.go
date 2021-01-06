package models

import (
	"time"
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
