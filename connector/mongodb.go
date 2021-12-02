package connector

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoInstance contains the Mongo client and database objects
type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var DB MongoInstance

// ConnectMongoDB configures the MongoDB client and initializes the database connection.
func ConnectMongoDB(host, dbName string) error {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + host + ":27017/" + dbName))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	db := client.Database(dbName)

	if err != nil {
		return err
	}

	DB = MongoInstance{
		Client: client,
		Db:     db,
	}

	return nil
}

func CreateNode(node *Node) error {
	coll := DB.Db.Collection("nodes")
	node.ID = ""

	_, err := coll.InsertOne(context.TODO(), node)
	if err != nil {
		return err
	}
	return nil

}

func FindNodeByWorkerID(worker_id string, result *Node) error {
	col := DB.Db.Collection("nodes")
	query := bson.D{{Key: "worker_id", Value: worker_id}}

	if err := col.FindOne(context.TODO(), query).Decode(result); err != nil {
		return err
	}
	return nil
}

func GetNodes(c *fiber.Ctx) error {
	query := bson.D{{}}
	cursor, err := DB.Db.Collection("nodes").Find(c.Context(), query)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	var nodes []Node = make([]Node, 0)
	if err := cursor.All(c.Context(), &nodes); err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.JSON(nodes)
}
