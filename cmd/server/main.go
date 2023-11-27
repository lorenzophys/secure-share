package main

import (
	"flag"
	"log"

	"github.com/lorenzophys/secure_share/internal/store"
	memoryStore "github.com/lorenzophys/secure_share/internal/store/in-memory"
)

type Config struct {
	Port      int
	StoreType string
	Debug     bool
	BaseUrl   string
}

type Application struct {
	Config Config
	Store  store.SecretStore
}

func main() {
	var cfg Config

	flag.IntVar(&cfg.Port, "port", 8080, "HTTP port")
	flag.BoolVar(&cfg.Debug, "debug", false, "Debug mode")
	flag.StringVar(&cfg.BaseUrl, "base-url", "http://localhost:8080", "Base URL")
	flag.StringVar(&cfg.StoreType, "store-type", "in-memory", "Secret store type")

	flag.Parse()

	var store store.SecretStore

	switch cfg.StoreType {
	case "in-memory":
		store = memoryStore.NewMemoryStore()
	default:
		log.Fatal("Invalid storage type")
	}

	app := &Application{
		Config: cfg,
		Store:  store,
	}

	err := app.serve()
	if err != nil {
		log.Fatal(err)
	}
}
