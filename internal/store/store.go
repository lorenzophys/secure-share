package store

import (
	"crypto/sha256"
	"time"

	"github.com/fernet/fernet-go"
)

type SecretStore interface {
	Set(value string, ttl time.Duration) string
	Get(key string) (string, bool)
	RemoveExpiredSecrets()
}

func GenerateFernetKeyFromUUID(uuid string) *fernet.Key {
	hash := sha256.Sum256([]byte(uuid))
	var key fernet.Key
	copy(key[:], hash[:])
	return &key
}
