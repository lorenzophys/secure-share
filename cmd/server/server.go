package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type SlogAdapter struct {
	*slog.Logger
}

func (s *SlogAdapter) Write(p []byte) (n int, err error) {
	s.Logger.Error(string(p))
	return len(p), nil
}

func (app *Application) NewTLSConfig() *tls.Config {
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	if cert, err := tls.LoadX509KeyPair(app.Config.TLS.CertFile, app.Config.TLS.KeyFile); err != nil {
		app.logger.Error("tls config error: failed to load key pair: %v", err)
	} else {
		tlsConfig.Certificates = []tls.Certificate{cert}
		app.logger.Info("tls config setup correctly.", "certFile", app.Config.TLS.CertFile, "keyFile", app.Config.TLS.KeyFile)
	}

	return tlsConfig
}

func (app *Application) Serve() error {
	var tlsConfig *tls.Config

	if app.Config.TLS.Enabled {
		tlsConfig = app.NewTLSConfig()
	}

	logAdapter := &SlogAdapter{app.logger}

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s", app.Config.ServicePort),
		Handler:      app.NewRouter(templatesPath),
		TLSConfig:    tlsConfig,
		ErrorLog:     log.New(logAdapter, "", 0),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		sig := <-quit

		app.logger.Info("shutting down the server.", "signal", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			shutdownError <- err
		}

		shutdownError <- nil
	}()

	var err error
	if app.Config.TLS.Enabled {
		app.logger.Info("starting TLS server...", "port", app.Config.ServicePort)
		err = srv.ListenAndServeTLS("", "")
	} else {
		app.logger.Info("starting server...", "port", app.Config.ServicePort)
		err = srv.ListenAndServe()
	}
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Info("stopped server.", "port", app.Config.ServicePort)

	return nil
}
