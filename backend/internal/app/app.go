// Package app wires the service together and runs it.
package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/ozaitsev92/tododdd/config"
	v1 "github.com/ozaitsev92/tododdd/internal/controller/http/v1"
	"github.com/ozaitsev92/tododdd/internal/controller/http/v1/jwt"
	"github.com/ozaitsev92/tododdd/internal/domain/task"
	"github.com/ozaitsev92/tododdd/internal/domain/user"
	"github.com/ozaitsev92/tododdd/internal/infrastructure/password"
	taskmemory "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/task/memory"
	taskmongo "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/task/mongo"
	usermemory "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/user/memory"
	usermongo "github.com/ozaitsev92/tododdd/internal/infrastructure/repository/user/mongo"
	"github.com/ozaitsev92/tododdd/internal/usecase"
	"github.com/ozaitsev92/tododdd/pkg/httpserver"
	"github.com/ozaitsev92/tododdd/pkg/mongodb"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Run serves until SIGINT/SIGTERM or a server error, then drains HTTP and
// closes the database.
func Run(cfg config.Config, log *slog.Logger) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	store, err := openStorage(ctx, cfg, log)
	if err != nil {
		return err
	}
	defer store.close(context.WithoutCancel(ctx), log)

	hasher, err := password.NewBcryptHasher(cfg.PasswordHashCost)
	if err != nil {
		return err
	}

	users, err := usecase.NewUserUseCase(store.users, hasher)
	if err != nil {
		return err
	}

	tokens, err := jwt.NewService(jwt.Config{
		SigningKey:   []byte(cfg.JWTSigningKey),
		Issuer:       cfg.JWTIssuer,
		TTL:          cfg.JWTSession(),
		CookieDomain: cfg.JWTCookieDomain,
		SecureCookie: cfg.JWTSecureCookie,
	})
	if err != nil {
		return err
	}

	gin.SetMode(gin.ReleaseMode)
	handler := v1.NewRouter(v1.Deps{
		Log:          log,
		Tokens:       tokens,
		Tasks:        usecase.NewTaskUseCase(store.tasks),
		Users:        users,
		MaxBodyBytes: cfg.MaxBodyBytes,
		Ready:        store.ready,
	})

	server := httpserver.New(handler,
		httpserver.Port(cfg.BindAddr),
		httpserver.ReadHeaderTimeout(config.Seconds(cfg.ReadHeaderTimeout)),
		httpserver.ReadTimeout(config.Seconds(cfg.ReadTimeout)),
		httpserver.WriteTimeout(config.Seconds(cfg.WriteTimeout)),
		httpserver.IdleTimeout(config.Seconds(cfg.IdleTimeout)),
		httpserver.ShutdownTimeout(config.Seconds(cfg.GracefulTimeout)),
	)
	server.Start()
	log.Info("http server started", "addr", server.Addr(), "storage", cfg.Storage)

	var serveErr error
	select {
	case <-ctx.Done():
		log.Info("shutdown signal received")
	case err, ok := <-server.Notify():
		if ok {
			serveErr = fmt.Errorf("http server: %w", err)
		}
	}

	if err := server.Shutdown(); err != nil {
		serveErr = errors.Join(serveErr, fmt.Errorf("http server shutdown: %w", err))
	}
	log.Info("http server stopped")

	return serveErr
}

type storage struct {
	tasks task.Repository
	users user.Repository
	ready func(context.Context) error
	close func(context.Context, *slog.Logger)
}

func openStorage(ctx context.Context, cfg config.Config, log *slog.Logger) (storage, error) {
	if cfg.Storage == config.StorageMemory {
		log.Warn("using in-memory storage: data is lost on restart")
		return storage{
			tasks: taskmemory.NewRepository(),
			users: usermemory.NewRepository(),
			close: func(context.Context, *slog.Logger) {},
		}, nil
	}

	timeout := config.Seconds(cfg.MongoTimeout)

	client, err := mongodb.Connect(ctx, cfg.MongoURL, timeout)
	if err != nil {
		return storage{}, err
	}

	closeClient := func(ctx context.Context, log *slog.Logger) {
		dctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		if err := client.Disconnect(dctx); err != nil {
			log.Error("mongodb disconnect", "err", err)
			return
		}
		log.Info("mongodb disconnected")
	}

	db := client.Database(cfg.MongoDBName)
	tasks := taskmongo.NewRepository(db)
	users := usermongo.NewRepository(db)

	if err := errors.Join(tasks.EnsureIndexes(ctx), users.EnsureIndexes(ctx)); err != nil {
		closeClient(context.WithoutCancel(ctx), log)
		return storage{}, err
	}

	return storage{
		tasks: tasks,
		users: users,
		ready: func(ctx context.Context) error { return client.Ping(ctx, readpref.Primary()) },
		close: closeClient,
	}, nil
}
