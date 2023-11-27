package redis_store

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/fernet/fernet-go"
	"github.com/google/uuid"
	"github.com/lorenzophys/secure_share/internal/store"
	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(addr, password string, db int) *RedisStore {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisStore{client: rdb}
}

func (rs *RedisStore) Set(value string, ttl time.Duration) string {
	ctx := context.Background()

	urlKey := uuid.New().String()

	fernetKey := store.GenerateFernetKeyFromUUID(urlKey)

	encryptedSecret, err := fernet.EncryptAndSign([]byte(value), fernetKey)
	if err != nil {
		log.Fatal(err)
	}

	truncatedKey := strings.Split(urlKey, "-")[0]

	err = rs.client.Set(ctx, truncatedKey, encryptedSecret, ttl).Err()
	if err != nil {
		log.Fatal(err)
	}

	return urlKey
}

func (rs *RedisStore) Get(urlKey string) (string, bool) {
	ctx := context.Background()
	truncatedKey := strings.Split(urlKey, "-")[0]

	encryptedSecret, err := rs.client.Get(ctx, truncatedKey).Result()
	if err != nil {
		return "", false
	}

	fernetKey := store.GenerateFernetKeyFromUUID(urlKey)
	fernetKeyList := []*fernet.Key{fernetKey}

	secret := fernet.VerifyAndDecrypt([]byte(encryptedSecret), 0, fernetKeyList)
	if secret == nil {
		return "", false
	}

	err = rs.client.Del(ctx, truncatedKey).Err()
	if err != nil {
		log.Fatal(err)
	}

	return string(secret), true
}

func (rs *RedisStore) RemoveExpiredSecrets() {}
