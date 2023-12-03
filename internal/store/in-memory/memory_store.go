package memory_store

import (
	"crypto/rand"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fernet/fernet-go"
	"github.com/google/uuid"
	"github.com/lorenzophys/secure_share/internal/store"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

type secretItem struct {
	secret    string
	ttl       time.Duration
	timeStamp time.Time
}

type MemoryStore struct {
	store map[string]secretItem
	sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	ms := &MemoryStore{
		store: make(map[string]secretItem),
	}

	go ms.RemoveExpiredSecrets()
	logger.Info("new in-memory store created successfully.")

	return ms
}

func (ms *MemoryStore) RemoveExpiredSecrets() {
	ticker := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-ticker.C:
			ms.Lock()
			defer ms.Unlock()

			for key, item := range ms.store {
				if time.Since(item.timeStamp) > item.ttl {
					delete(ms.store, key)
					logger.Info("deleted expired key", "expired_key", key)
				}
			}
		}
	}

}

func (ps *MemoryStore) Set(value string, ttl time.Duration) string {
	ps.Lock()
	defer ps.Unlock()

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
	ps.store[truncatedKey] = secretItem{
		secret:    string(encryptedSecretWithSalt),
		ttl:       ttl,
		timeStamp: time.Now(),
	}

	return urlKey
}

func (ps *MemoryStore) Get(urlKey string) (string, bool) {
	ps.RLock()
	defer ps.RUnlock()

	truncatedKey := strings.Split(urlKey, "-")[0]
	encryptedSecretWithSalt, ok := ps.store[truncatedKey]
	if !ok {
		return "", false
	}

	salt := encryptedSecretWithSalt.secret[:16]
	encryptedSecret := encryptedSecretWithSalt.secret[16:]

	fernetKey := store.GenerateFernetKeyFromUUID(urlKey, []byte(salt))

	fernetKeyList := []*fernet.Key{fernetKey}
	secret := fernet.VerifyAndDecrypt([]byte(encryptedSecret), encryptedSecretWithSalt.ttl, fernetKeyList)
	if secret == nil {
		return "", false
	}

	delete(ps.store, truncatedKey)
	logger.Info("secret revealed, hence deleted from the store", "secret_key", urlKey)

	return string(secret), true
}
