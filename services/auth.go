package services

import (
	"crypto/rsa"
	"fmt"
	"github.com/github.com/steevehook/account-api/models"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/lestrrat-go/jwx/jwt"
	"log"
)

// AccountsRepository represents the accounts repository
type AccountsRepository interface {
	GetPrivateKey() (*rsa.PrivateKey, error)
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
func (s Auth) Signup(credentials models.Credentials) (models.TokenResponse, error) {
	// create new private key with expiration
	// rotate the private key creation every 2h
	// save the jwk new and old keys in Redis
	// create new key when old one is expired
	key, err := s.repo.GetPrivateKey()
	if err != nil {
		return models.TokenResponse{}, err
	}
	fmt.Println(key)

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
