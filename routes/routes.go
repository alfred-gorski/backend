package routes

import (
	"backend/controller"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	nodes := api.Group("/nodes")
	nodes.Get("/", controller.GetNodes)
	nodes.Get("/:id", controller.GetNode)
	nodes.Put("/:id", controller.UpdateNode)
}
