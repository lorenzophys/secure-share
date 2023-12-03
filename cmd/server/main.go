package main

import (
	"log"
	"log/slog"
	"os"
	"strconv"

	"github.com/lorenzophys/secure_share/internal/store"
	memoryStore "github.com/lorenzophys/secure_share/internal/store/in-memory"
	redisStore "github.com/lorenzophys/secure_share/internal/store/redis"
)

type Application struct {
	Config Config
	Store  store.SecretStore
	logger *slog.Logger
}

func main() {
	var store store.SecretStore

	app := &Application{
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}

	app.Config = NewConfig()

	switch app.Config.StoreBackend {
	case "in-memory":
		store = memoryStore.NewMemoryStore()
	case "redis":
		var value int64
		value, err := strconv.ParseInt(app.Config.Redis.Db, 10, 64)
		if err != nil {
			app.logger.Error("error parsing store type as int. Using default value: 0", "redis DB", app.Config.Redis.Db, "error", err)
		}
		store = redisStore.NewRedisStore(app.Config.Redis.Address, app.Config.Redis.Password, int(value))
	default:
		app.logger.Error("invalid storage type", "storage_type", app.Config.StoreBackend)
		os.Exit(1)
	}

	app.Store = store

	err := app.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
