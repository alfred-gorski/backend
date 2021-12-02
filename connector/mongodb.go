package connector

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB contains the Mongo database objects
var DB *mongo.Database
var client *mongo.Client

// ConnectMongoDB configures the MongoDB client and initializes the database connection.
func ConnectMongoDB(host, dbName string) error {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + host + ":27017/" + dbName))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	DB = client.Database(dbName)

	if err != nil {
		return err
	}

	return nil
}

func UpdateNode(c *fiber.Ctx) error {
	idParam := c.Params("id")
	nodeID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.SendStatus(400)
	}
	node := new(Node)
	if err := c.BodyParser(node); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	query := bson.D{{Key: "_id", Value: nodeID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "schema", Value: node.Schema}}}}

	if err := DB.Collection("nodes").FindOneAndUpdate(c.Context(), query, update).Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return c.SendStatus(404)
		}
		return c.SendStatus(500)
	}
	node.ID = idParam
	return c.JSON(node)
}

func GetNodes(c *fiber.Ctx) error {
	query := bson.D{{}}
	cursor, err := DB.Collection("nodes").Find(c.Context(), query)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	var nodes = make([]Node, 0)
	if err := cursor.All(c.Context(), &nodes); err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.JSON(nodes)
}

func FindNodeByWorkerID(worker_id string, result *Node) error {
	col := DB.Collection("nodes")
	query := bson.D{{Key: "worker_id", Value: worker_id}}

	if err := col.FindOne(context.TODO(), query).Decode(result); err != nil {
		return err
	}
	return nil
}

func createNode(node *Node) error {
	coll := DB.Collection("nodes")
	node.ID = ""

	_, err := coll.InsertOne(context.TODO(), node)
	if err != nil {
		return err
	}
	return nil

}
