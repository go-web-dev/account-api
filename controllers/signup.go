package controllers

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/lestrrat-go/jwx/jwt"

	"github.com/github.com/steevehook/account-api/models"
)

type signupper interface {
	Signup(credentials models.Credentials)
}

func signup(service signupper) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create new private key with expiration
		// rotate the private key creation every 2h
		// save the jwk new and old keys in Redis
		// create new key when old one is expired

		token := jwt.New()
		key, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			log.Printf("failed to generate private key: %s", err)
			return
		}

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
			return
		}
		err = jwk.AssignKeyID(jwkKey, jws.WithHeaders(jws.NewHeaders()))
		if err != nil {
			log.Fatal(err)
		}

		s := jwk.Set{Keys: []jwk.Key{jwkKey}}

		jsonbuf, err := json.MarshalIndent(s, "", "  ")
		if err != nil {
			log.Printf("failed to generate JSON: %s", err)
			return
		}

		os.Stdout.Write(jsonbuf)

		signed, err := jwt.Sign(token, jwa.RS256, jwkKey)
		if err != nil {
			log.Printf("failed to sign token: %s", err)
			return
		}
		//fmt.Println("kid", jwkKey.KeyID())
		//fmt.Println("kid", jwkKey.PrivateParams())
		fmt.Println("jwk", string(signed))
	})
}
