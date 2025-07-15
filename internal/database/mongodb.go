package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type MongoDB struct {
	uri    string
	client *mongo.Client
	ctx    context.Context
	cancel context.CancelFunc
}

func NewMongoDB(uri string) *MongoDB {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	return &MongoDB{uri: uri, ctx: ctx, cancel: cancel}
}

func (m *MongoDB) Disconnect() error {
	if err := m.client.Disconnect(m.ctx); err != nil {
		return fmt.Errorf("fail on disconnect mongodb database: %v", err)
	}
	return nil
}

func (m *MongoDB) Connect() (*mongo.Client, error) {
	options := options.Client().ApplyURI(m.uri)

	client, err := mongo.Connect(options)
	if err != nil {
		return nil, fmt.Errorf("fail on connect to mongodb: %v", err)
	}

	if err := client.Ping(m.ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("fail on connect database: %v", err)
	}

	log.Println("mongo database has connected")

	m.client = client

	return client, err
}
