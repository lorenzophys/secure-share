package memory_store

import (
	"log"
	"strings"
	"sync"

	"github.com/fernet/fernet-go"
	"github.com/google/uuid"
	"github.com/lorenzophys/secure_share/internal/store"
)

type MemoryStore struct {
	store map[string]string
	sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		store: make(map[string]string),
	}
}

func (ps *MemoryStore) Set(value string) string {
	ps.Lock()
	defer ps.Unlock()

	urlKey := uuid.New().String()

	fernetKey := store.GenerateFernetKeyFromUUID(urlKey)

	encryptedSecret, err := fernet.EncryptAndSign([]byte(value), fernetKey)
	if err != nil {
		log.Fatal(err)
	}

	truncatedKey := strings.Split(urlKey, "-")[0]
	ps.store[truncatedKey] = string(encryptedSecret)

	return urlKey
}

func (ps *MemoryStore) Get(urlKey string) (string, bool) {
	ps.RLock()
	defer ps.RUnlock()

	truncatedKey := strings.Split(urlKey, "-")[0]
	encryptedSecret, ok := ps.store[truncatedKey]
	if !ok {
		return "", false
	}

	fernetKey := store.GenerateFernetKeyFromUUID(urlKey)
	fernetKeyList := []*fernet.Key{fernetKey}
	secret := fernet.VerifyAndDecrypt([]byte(encryptedSecret), 0, fernetKeyList)
	if secret == nil {
		log.Fatal("secret is nil")
	}

	delete(ps.store, truncatedKey)

	return string(secret), true
}
