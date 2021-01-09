package services

import (
	"crypto/rsa"
	"fmt"
	"log"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/lestrrat-go/jwx/jwt"

	"github.com/github.com/steevehook/account-api/models"
)

// AccountsRepository represents the Accounts repository
type AccountsRepository interface {
}

// KeysRepository represents the Keys repository
type KeysRepository interface {
	GetPrivateKey() (*rsa.PrivateKey, error)
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
func (s Auth) Signup(credentials models.Credentials) (models.TokenResponse, error) {
	// create new private key with expiration
	// rotate the private key creation every 2h
	// save the jwk new and old keys in Redis
	// create new key when old one is expired

	// generate new private key on each app startup
	// save the old key and the new key and only sign with the new key in Redis
	// save all public keys under 'keys' aka *jwk.Set in Redis
	// create a middleware library

	// add client_id, client_secret, grant_type and refresh_token possibilities
	// add refresh token rotation
	// add refresh token throttling

	key, err := s.keysRepo.GetPrivateKey()
	if err != nil {
		return models.TokenResponse{}, err
	}

	token := jwt.New()

	//bs, _ := json.Marshal(token)
	//signed, err := jws.Sign(bs, jwa.RS256, key, jws.WithHeaders(jws.NewHeaders()))
	//if err != nil {
	//	log.Printf("failed to sign token: %s", err)
	//	return
	//}
	//fmt.Println(jws.NewHeaders())
	//fmt.Println("key:", string(signed))

	//privKey := ExportRsaPrivateKeyAsPemStr(key)
	//fmt.Println("priv:", privKey)
	//
	//pub := rsa.PublicKey{
	//	N: key.N,
	//	E: key.E,
	//}
	//pubKey, _ := ExportRsaPublicKeyAsPemStr(&pub)
	//fmt.Println("pub:", pubKey)

	//signed, err := jwt.Sign(token, jwa.RS256, key)
	//if err != nil {
	//	log.Printf("failed to sign token: %s", err)
	//	return
	//}
	//fmt.Println("jwt", string(signed))

	jwkKey, err := jwk.New(key)
	if err != nil {
		log.Printf("failed to create JWK key: %s", err)
	}
	err = jwk.AssignKeyID(jwkKey, jws.WithHeaders(jws.NewHeaders()))
	if err != nil {
		log.Fatal(err)
	}
	//
	//ss := jwk.Set{Keys: []jwk.Key{jwkKey}}
	//jsonbuf, err := json.MarshalIndent(ss, "", "  ")
	//if err != nil {
	//	log.Printf("failed to generate JSON: %s", err)
	//}
	//
	//os.Stdout.Write(jsonbuf)
	//
	signed, err := jwt.Sign(token, jwa.RS256, jwkKey)
	if err != nil {
		log.Printf("failed to sign token: %s", err)
	}
	//fmt.Println("kid", jwkKey.KeyID())
	//fmt.Println("kid", jwkKey.PrivateParams())
	fmt.Println("jwk", string(signed))

	return models.TokenResponse{}, nil
}

// Logout logs out a specific user account
func (s Auth) Logout() {
}
