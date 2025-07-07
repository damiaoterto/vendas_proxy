package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func ConnectFromURI(uri string) (*mongo.Client, error) {
	options := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(options)
	if err != nil {
		return nil, fmt.Errorf("fail on connect to mongodb: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_ = client.Ping(ctx, readpref.Primary())

	defer client.Disconnect(ctx)

	return client, err
}
