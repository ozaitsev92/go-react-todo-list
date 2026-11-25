// Package mongodb creates MongoDB clients.
package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Connect returns a client after pinging the server; callers must Disconnect it.
func Connect(ctx context.Context, uri string, timeout time.Duration) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI(uri).
		SetConnectTimeout(timeout).
		SetServerSelectionTimeout(timeout).
		SetTimeout(timeout)) // default per-operation timeout (CSOT)
	if err != nil {
		return nil, fmt.Errorf("mongodb: connect: %w", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		dctx, dcancel := context.WithTimeout(context.WithoutCancel(ctx), timeout)
		defer dcancel()
		_ = client.Disconnect(dctx)
		return nil, fmt.Errorf("mongodb: ping: %w", err)
	}

	return client, nil
}
