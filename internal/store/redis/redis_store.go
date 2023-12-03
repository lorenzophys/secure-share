package redis_store

import (
	"context"
	"crypto/rand"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/fernet/fernet-go"
	"github.com/google/uuid"
	"github.com/lorenzophys/secure_share/internal/store"
	"github.com/redis/go-redis/v9"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(addr, password string, db int) *RedisStore {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	logger.Info("new Redis store created successfully.")

	return &RedisStore{client: rdb}
}

func (rs *RedisStore) Set(value string, ttl time.Duration) string {
	ctx := context.Background()

	urlKey := uuid.New().String()

	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		logger.Error("failed to generate random salt.", "salt", salt, "error", err)
	}

	fernetKey := store.GenerateFernetKeyFromUUID(urlKey, salt)

	encryptedSecret, err := fernet.EncryptAndSign([]byte(value), fernetKey)
	if err != nil {
		logger.Error("failed to encrypt secret", "error", err)
	}

	encryptedSecretWithSalt := append(salt, encryptedSecret...)

	truncatedKey := strings.Split(urlKey, "-")[0]

	err = rs.client.Set(ctx, truncatedKey, encryptedSecretWithSalt, ttl).Err()
	if err != nil {
		logger.Error("call to Redis SET failed", "error", err)
	}

	return urlKey
}

func (rs *RedisStore) Get(urlKey string) (string, bool) {
	ctx := context.Background()
	truncatedKey := strings.Split(urlKey, "-")[0]

	encryptedSecretWithSalt, err := rs.client.Get(ctx, truncatedKey).Result()
	if err != nil {
		return "", false
	}

	salt := encryptedSecretWithSalt[:16]
	encryptedSecret := encryptedSecretWithSalt[16:]

	fernetKey := store.GenerateFernetKeyFromUUID(urlKey, []byte(salt))
	fernetKeyList := []*fernet.Key{fernetKey}

	secret := fernet.VerifyAndDecrypt([]byte(encryptedSecret), 0, fernetKeyList)
	if secret == nil {
		return "", false
	}

	err = rs.client.Del(ctx, truncatedKey).Err()
	if err != nil {
		logger.Error("call to Redis DEL failed", "error", err)
	}
	logger.Info("secret revealed, hence deleted from the store", "secret_key", urlKey)

	return string(secret), true
}

func (rs *RedisStore) RemoveExpiredSecrets() {}
