package mymongo

import (
	"StudShare/internal/config/storage_config"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func NewMongo(cfg storage_config.MongoConfig) (*mongo.Database, error) {
	clientOpts := options.Client().ApplyURI(cfg.URI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client.Database(cfg.DBName), nil
}
