package services

import (
	"context"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"go.uber.org/zap"

	"github.com/github.com/steevehook/account-api/logging"
	"github.com/github.com/steevehook/account-api/models"
)

// AccountsRepository represents the Accounts repository
type AccountsRepository interface {
}

// KeysRepository represents the Keys repository
type KeysRepository interface {
	GetPrivateKey(ctx context.Context) (jwk.RSAPrivateKey, error)
	GetKeySet(ctx context.Context) (jwk.Set, error)
}

// NewAuth creates a new instance of Auth service
func NewAuth(accountsRepo AccountsRepository, keysRepo KeysRepository) Auth {
	service := Auth{
		accountsRepo: accountsRepo,
		keysRepo:     keysRepo,
	}
	return service
}

// Auth represents the authentication service
type Auth struct {
	accountsRepo AccountsRepository
	keysRepo     KeysRepository
}

// Login logins the given user
func (s Auth) Login() {
}

// Signup creates a new user account
func (s Auth) Signup(ctx context.Context, credentials models.Credentials) (models.TokenResponse, error) {
	// create a middleware library

	// add client_id, client_secret, grant_type and refresh_token possibilities
	// add refresh token rotation
	// add refresh token throttling

	logger := logging.Logger
	token := models.NewJWT()
	key, err := s.keysRepo.GetPrivateKey(ctx)
	if err != nil {
		return models.TokenResponse{}, err
	}

	signed, err := jwt.Sign(token, jwa.RS256, key)
	if err != nil {
		logger.Error("could not sign jwt token", zap.Error(err))
		return models.TokenResponse{}, err
	}

	res := models.TokenResponse{
		AccessToken: string(signed),
	}
	return res, nil
}

// Logout logs out a specific user account
func (s Auth) Logout() {
}

func (s Auth) GetKeySet(ctx context.Context) (jwk.Set, error) {
	return s.keysRepo.GetKeySet(ctx)
}
