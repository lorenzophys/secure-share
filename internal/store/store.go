package store

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/fernet/fernet-go"
	"golang.org/x/crypto/pbkdf2"
)

type SecretStore interface {
	Set(value string, ttl time.Duration) string
	Get(key string) (string, bool)
	RemoveExpiredSecrets()
}

func GenerateFernetKeyFromUUID(uuid string, salt []byte) *fernet.Key {
	iterations := 600000 // https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html#pbkdf2
	keyLength := 32      // 256 bits for Fernet key
	pbkdf2Key := pbkdf2.Key([]byte(uuid), salt, iterations, keyLength, sha256.New)
	pbkdf2KeyHex := hex.EncodeToString(pbkdf2Key)

	var key fernet.Key
	copy(key[:], pbkdf2KeyHex[:])

	return &key
}
