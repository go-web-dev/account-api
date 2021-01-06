package services

import (
	"github.com/github.com/steevehook/account-api/models"
)

// AccountsRepository represents the accounts repository
type AccountsRepository interface {
}

// NewAuth creates a new instance of Auth service
func NewAuth(repo AccountsRepository) Auth {
	service := Auth{
		repo: repo,
	}
	return service
}

// Auth represents the authentication service
type Auth struct {
	repo AccountsRepository
}

// Login logins the given user
func (s Auth) Login() {
}

// Signup creates a new user account
func (s Auth) Signup(credentials models.Credentials) {
}

// Logout logs out a specific user account
func (s Auth) Logout() {
}
