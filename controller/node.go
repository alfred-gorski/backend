package controller

import (
	mongodb_conn "backend/connector/mongodb"
	"backend/models"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateNode(c *fiber.Ctx) error {
	idParam := c.Params("id")
	nodeID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.SendStatus(400)
	}
	node := new(models.Node)
	if err := c.BodyParser(node); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	query := bson.D{{Key: "_id", Value: nodeID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "schema", Value: node.Schema}, {Key: "form_data", Value: node.FormData}}}}

	if err := mongodb_conn.DB.Collection("nodes").FindOneAndUpdate(c.Context(), query, update).Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return c.SendStatus(404)
		}
		return c.SendStatus(500)
	}
	node.ID = idParam
	return c.JSON(node)
}

func GetNode(c *fiber.Ctx) error {
	idParam := c.Params("id")
	nodeID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.SendStatus(400)
	}
	query := bson.D{{Key: "_id", Value: nodeID}}
	var result = models.Node{}
	err = mongodb_conn.DB.Collection("nodes").FindOne(c.Context(), query).Decode(&result)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.JSON(result)
}

func GetNodes(c *fiber.Ctx) error {
	query := bson.D{{}}
	cursor, err := mongodb_conn.DB.Collection("nodes").Find(c.Context(), query)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	var nodes = make([]models.Node, 0)
	if err := cursor.All(c.Context(), &nodes); err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.JSON(nodes)
}

func FindNodeByWorkerID(worker_id string, result *models.Node) error {
	col := mongodb_conn.DB.Collection("nodes")
	query := bson.D{{Key: "worker_id", Value: worker_id}}

	if err := col.FindOne(context.TODO(), query).Decode(result); err != nil {
		return err
	}
	return nil
}

func CreateNode(node *models.Node) error {
	coll := mongodb_conn.DB.Collection("nodes")
	node.ID = ""

	_, err := coll.InsertOne(context.TODO(), node)
	if err != nil {
		return err
	}
	return nil

}
