package routes

import (
	"backend/connector"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	nodes := api.Group("/nodes")
	nodes.Get("/", connector.GetNodes)
	nodes.Put("/:id", connector.UpdateNode)
}
