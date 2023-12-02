package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func (app *Application) NewTLSConfig() *tls.Config {
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	if cert, err := tls.LoadX509KeyPair(app.Config.TLS.CertFile, app.Config.TLS.KeyFile); err != nil {
		log.Printf("Failed to load key pair: %v", err)
	} else {
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	return tlsConfig
}

func (app *Application) Serve() error {
	var tlsConfig *tls.Config

	if app.Config.TLS.Enabled {
		tlsConfig = app.NewTLSConfig()
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s", app.Config.ServicePort),
		Handler:      app.NewRouter(templatesPath),
		TLSConfig:    tlsConfig,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		sig := <-quit

		log.Printf("shutting down the server: received %s", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			shutdownError <- err
		}

		shutdownError <- nil
	}()

	log.Printf("starting server on port %s", app.Config.ServicePort)

	var err error
	if app.Config.TLS.Enabled {
		err = srv.ListenAndServeTLS("", "")
	} else {
		err = srv.ListenAndServe()
	}
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	log.Printf("stopped server on port %s", app.Config.ServicePort)

	return nil
}
