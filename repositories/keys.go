package repositories

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jws"
	"go.uber.org/zap"

	"github.com/github.com/steevehook/account-api/logging"
)

const (
	redisPrivateKeyKey = "private_key"
	redisPrivateKeyExp = time.Hour
	redisKeySetKey     = "keys"
	redisKeySetCap     = 20
)

// CacheDriver represents the application cache driver
type CacheDriver interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	LTrim(ctx context.Context, key string, start, stop int64) *redis.StatusCmd
	LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd
	Close() error
}

// Keys represents the keys repository
type Keys struct {
	cache CacheDriver
}

// NewKeys creates a new Keys repository
func NewKeys(cache CacheDriver) Keys {
	repo := Keys{
		cache: cache,
	}
	return repo
}

// GetPrivateKey fetches the existing private key or creates a new one
func (repo Keys) GetPrivateKey(ctx context.Context) (jwk.RSAPrivateKey, error) {
	logger := logging.Logger
	privateJWKBytes, err := repo.cache.Get(ctx, redisPrivateKeyKey).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if len(privateJWKBytes) != 0 {
		key := jwk.NewRSAPrivateKey()
		err = json.Unmarshal(privateJWKBytes, &key)
		if err != nil {
			return nil, err
		}
		return key, nil
	}

	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		logger.Error("failed to generate rsa key", zap.Error(err))
		return nil, err
	}
	privateJWK, publicJWK, err := repo.newJWKRSAKeys(rsaPrivateKey)
	if err != nil {
		logger.Error("could not create rsa jwk keys from rsa private key", zap.Error(err))
		return nil, err
	}

	err = repo.savePrivateJWK(ctx, privateJWK)
	if err != nil {
		logger.Error("could not save new private jwk", zap.Error(err))
		return nil, err
	}
	err = repo.savePublicJWK(ctx, publicJWK)
	if err != nil {
		logger.Error("could not save new public jwk", zap.Error(err))
		return nil, err
	}

	return privateJWK, nil
}

// GetKeySet fetches public jwk key set
func (repo Keys) GetKeySet(ctx context.Context) (jwk.Set, error) {
	logger := logging.Logger
	res, err := repo.cache.LRange(ctx, redisKeySetKey, 0, redisKeySetCap).Result()
	if err != nil {
		logger.Error("could not fetch public jwk key set", zap.Error(err))
		return jwk.Set{}, err
	}

	var keys []jwk.Key
	for _, k := range res {
		key, err := jwk.ParseKey([]byte(k))
		if err != nil {
			logger.Error("could not parse jwk key", zap.Error(err))
			return jwk.Set{}, err
		}
		keys = append(keys, key)
	}
	return jwk.Set{Keys: keys}, nil
}

func (repo Keys) newJWKRSAKeys(privateKey *rsa.PrivateKey) (jwk.RSAPrivateKey, jwk.RSAPublicKey, error) {
	logger := logging.Logger
	privateJWK := jwk.NewRSAPrivateKey()
	err := repo.setHeaders(privateJWK)
	if err != nil {
		return nil, nil, err
	}
	err = privateJWK.FromRaw(privateKey)
	if err != nil {
		logger.Error("could not set jwk raw private rsa key for jwk", zap.Error(err))
		return nil, nil, err
	}
	publicJWK, err := privateJWK.PublicKey()
	if err != nil {
		logger.Error("could not get public jwk", zap.Error(err))
		return nil, nil, err
	}
	err = repo.setHeaders(publicJWK)
	if err != nil {
		return nil, nil, err
	}
	return privateJWK, publicJWK, nil
}

func (repo Keys) setHeaders(key jwk.Key) error {
	logger := logging.Logger
	err := jwk.AssignKeyID(key, jws.WithHeaders(jws.NewHeaders()))
	if err != nil {
		logger.Error("could not set jws headers for jwk key", zap.Error(err))
		return err
	}
	err = key.Set(jwk.AlgorithmKey, jwa.RS256)
	if err != nil {
		logger.Error("could not set alg field for jwk key", zap.Error(err))
		return err
	}
	return nil
}

func (repo Keys) savePrivateJWK(ctx context.Context, privateJWK jwk.RSAPrivateKey) error {
	logger := logging.Logger
	privateJWKBytes, err := json.Marshal(privateJWK)
	if err != nil {
		logger.Error("could not marshal new private jwk", zap.Error(err))
		return err
	}
	err = repo.cache.Set(ctx, redisPrivateKeyKey, privateJWKBytes, redisPrivateKeyExp).Err()
	if err != nil {
		logger.Error("could not set new rsa private key", zap.Error(err))
		return err
	}
	return nil
}

func (repo Keys) savePublicJWK(ctx context.Context, publicJWK jwk.RSAPublicKey) error {
	logger := logging.Logger
	publicJWKBytes, err := json.Marshal(publicJWK)
	if err != nil {
		logger.Error("could not marshal new public jwk", zap.Error(err))
		return err
	}
	err = repo.cache.LPush(ctx, redisKeySetKey, publicJWKBytes).Err()
	if err != nil {
		logger.Error("could not push new public jwk to key set list", zap.Error(err))
		return err
	}
	err = repo.cache.LTrim(ctx, redisKeySetKey, 0, redisKeySetCap).Err()
	if err != nil {
		logger.Error("could not trim public jwk key set list", zap.Error(err))
		return err
	}
	return nil
}
