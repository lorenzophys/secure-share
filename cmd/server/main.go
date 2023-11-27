package main

import (
	"log"
	"strconv"

	"github.com/lorenzophys/secure_share/internal/store"
	memoryStore "github.com/lorenzophys/secure_share/internal/store/in-memory"
	redisStore "github.com/lorenzophys/secure_share/internal/store/redis"
)

type Application struct {
	Config Config
	Store  store.SecretStore
}

func main() {
	var store store.SecretStore
	cfg := NewConfig()

	switch cfg.StoreBackend {
	case "in-memory":
		store = memoryStore.NewMemoryStore()
	case "redis":
		var value int64
		value, err := strconv.ParseInt(cfg.RedisDb, 10, 64)
		if err != nil {
			log.Printf("Error parsing %s as int: %s. Using default value: 0", cfg.RedisDb, err)
		}
		store = redisStore.NewRedisStore(cfg.RedisAddr, cfg.RedisPassword, int(value))
	default:
		log.Fatal("Invalid storage type")
	}

	app := &Application{
		Config: *cfg,
		Store:  store,
	}

	err := app.serve()
	if err != nil {
		log.Fatal(err)
	}
}
