package repositories

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/lestrrat-go/jwx/jwk"
	"go.uber.org/zap"

	"github.com/github.com/steevehook/account-api/logging"
)

const redisPrivateKey = "private_key"

// DBDriver represents the database driver
type DBDriver interface {
	Close() error
}

type CacheDriver interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Close() error
}

// Accounts represents the accounts repository
type Accounts struct {
	db    DBDriver
	cache CacheDriver
}

// NewAccounts creates a new Accounts repository
func NewAccounts(db DBDriver, cache CacheDriver) Accounts {
	repo := Accounts{
		db:    db,
		cache: cache,
	}
	return repo
}

func (repo Accounts) GetPrivateKey() (*rsa.PrivateKey, error) {
	logger := logging.Logger
	ctx := context.Background()
	bs, err := repo.cache.Get(ctx, redisPrivateKey).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if len(bs) != 0 {
		jwkKey, err := jwk.ParseKey(bs)
		if err != nil {
			return nil, err
		}

		b, _ := json.Marshal(jwkKey)
		fmt.Println("JWK", string(b))

		kk := map[string]string{}
		json.Unmarshal(b, &kk)
		fmt.Println("D BEFORE:", kk["d"])

		d, err := base64.RawURLEncoding.DecodeString(kk["d"])
		if err != nil {
			log.Fatal(err)
		}
		// use q and p for private key
		// use dp dq and qi for precomputed values
		kkk := rsa.PrivateKey{
			D: new(big.Int).SetBytes(d),
		}
		kkk.Precompute()
		fmt.Println("D:", kkk.D)

		return nil, nil
	}
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		logger.Error("failed to generate rsa key", zap.Error(err))
		return nil, err
	}

	jwkKey, err := jwk.New(key)
	if err != nil {
		logger.Error("failed to create jwk key", zap.Error(err))
		return nil, err
	}
	bs, err = json.Marshal(jwkKey)
	if err != nil {
		logger.Error("failed to marshal jwk key", zap.Error(err))
		return nil, err
	}

	status := repo.cache.Set(ctx, redisPrivateKey, bs, time.Hour)
	if status.Err() != nil {
		logger.Error("could not set new rsa private key", zap.Error(status.Err()))
		return nil, err
	}

	return key, nil
}
