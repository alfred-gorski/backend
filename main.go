package main

import (
	"backend/app/routes"
	"backend/connector"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(logger.New())

	if err := connector.ConnectMongoDB(os.Args[1], os.Args[2]); err != nil {
		log.Fatal(err)
	}

	if err := connector.ConnectMQTT(os.Args[1]); err != nil {
		log.Fatal(err)
	}

	routes.SetupRoutes(app)

	app.Static("/", "dist")
	app.Get("/*", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("dist/index.html")
	})

	log.Fatal(app.Listen(":3000"))
}
