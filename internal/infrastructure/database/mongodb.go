package database

import (
	"context"
	"fmt"

	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDb struct {
	options *options.ClientOptions
	client  *mongo.Client
	Db      *mongo.Database
}

func NewMongodb() (*MongoDb, error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", config.Config.Database.Username, config.Config.Database.Password, config.Config.Database.Host, config.Config.Database.Port)

	return &MongoDb{
		options: options.
			Client().
			ApplyURI(uri),
	}, nil
}

func (m *MongoDb) Connect() error {
	client, err := mongo.Connect(m.options)
	if err != nil {
		return err
	}

	m.client = client

	err = m.client.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}

	m.Db = client.Database(config.Config.Database.Name)

	return nil
}

func (m *MongoDb) Disconnect(ctx context.Context) {
	m.client.Disconnect(ctx)
}
