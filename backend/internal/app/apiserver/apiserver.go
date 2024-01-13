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

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func Start(config *Config) {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file:///migrations", "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil {
		log.Fatal(err)
	}

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
