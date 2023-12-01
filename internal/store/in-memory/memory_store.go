package memory_store

import (
	"crypto/rand"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/fernet/fernet-go"
	"github.com/google/uuid"
	"github.com/lorenzophys/secure_share/internal/store"
)

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
		log.Printf("Failed to generate random salt: %v", err)
	}

	fernetKey := store.GenerateFernetKeyFromUUID(urlKey, salt)

	encryptedSecret, err := fernet.EncryptAndSign([]byte(value), fernetKey)
	if err != nil {
		log.Fatal(err)
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

	return string(secret), true
}
