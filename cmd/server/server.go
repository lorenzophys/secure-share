package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func (app *Application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s", app.Config.ServicePort),
		Handler:      app.NewRouter(templatesPath),
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
	err := srv.ListenAndServe()
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
