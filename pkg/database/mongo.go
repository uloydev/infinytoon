package database

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	appctx "infinitoon.dev/infinitoon/pkg/context"
)

type MongoDB struct {
	appCtx *appctx.AppContext
	client *mongo.Client
	URI    string
}

func NewMongoDB(appCtx *appctx.AppContext, uri string) *MongoDB {
	db := &MongoDB{
		URI: uri,
	}

	appCtx.Set(appctx.DatabaseKey, db)
	return db
}

func (m *MongoDB) Connect() error {
	// Connect to MongoDB
	client, err := mongo.Connect(options.Client().ApplyURI(m.URI))
	if err != nil {
		return err
	}

	// Check the connection
	err = client.Ping(m.appCtx.Context(), nil)
	if err != nil {
		return err
	}

	m.client = client
	return nil
}

func (m *MongoDB) Disconnect() error {
	err := m.client.Disconnect(m.appCtx.Context())
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoDB) Client() *mongo.Client {
	return m.client
}
