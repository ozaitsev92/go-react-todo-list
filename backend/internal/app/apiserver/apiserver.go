package apiserver

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ozaitsev92/go-react-todo-list/internal/app/apiserver/jwt"
	"github.com/ozaitsev92/go-react-todo-list/internal/app/infrastructure/store/sqlstore"
)

func Start(config *Config) {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	store := sqlstore.New(db)
	jwtService := jwt.NewJWTService([]byte(config.JWTSigningKey), config.JWTSessionLength, config.JWTCookieDomain, config.JWTSecureCookie)
	appServer := newServer(store, jwtService, config)

	srv := &http.Server{
		Addr:         config.BindAddr,
		WriteTimeout: time.Duration(config.GracefulTimeout) * time.Second,
		ReadTimeout:  time.Duration(config.GracefulTimeout) * time.Second,
		IdleTimeout:  time.Duration(config.GracefulTimeout) * time.Second,
		Handler:      appServer,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			appServer.logger.Error(err)
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.GracefulTimeout)*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		appServer.logger.Error(err)
	}

	appServer.logger.Info("shutting down")
	os.Exit(0)
}

func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
