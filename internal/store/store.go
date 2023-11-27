package store

import (
	"crypto/sha256"

	"github.com/fernet/fernet-go"
)

type SecretStore interface {
	Set(value string) string
	Get(key string) (string, bool)
}

func GenerateFernetKeyFromUUID(uuid string) *fernet.Key {
	hash := sha256.Sum256([]byte(uuid))
	var key fernet.Key
	copy(key[:], hash[:])
	return &key
}
