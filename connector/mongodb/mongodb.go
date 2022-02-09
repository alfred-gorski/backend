package mongodb_conn

import (
	"backend/config"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB contains the Mongo database objects
var DB *mongo.Database

// ConnectMongoDB configures the MongoDB client and initializes the database connection.
func ConnectMongoDB(m config.MongoConfig) error {
	uri := fmt.Sprintf("mongodb://%s:%d/%s", m.Addr, m.Port, m.DB)
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	DB = client.Database(m.DB)

	if err != nil {
		return err
	}

	return nil
}
