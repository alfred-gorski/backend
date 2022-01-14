package main

import (
	"backend/config"
	"backend/routes"
	"log"

	"github.com/gofiber/fiber/v2"

	mongodb_conn "backend/connector/mongodb"
	mqtt_conn "backend/connector/mqtt"
)

func main() {
	app := fiber.New()
	// app.Use(cors.New())
	// app.Use(recover.New())
	// app.Use(logger.New())
	var conf config.Config
	config.LoadConfig(&conf)

	if err := mongodb_conn.ConnectMongoDB(conf.Mongo); err != nil {
		log.Fatal(err)
	}

	if err := mqtt_conn.ConnectMQTT(conf.MQTT); err != nil {
		log.Fatal(err)
	}

	routes.SetupRoutes(app)

	app.Static("/", "dist")
	app.Get("/*", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("dist/index.html")
	})

	log.Fatal(app.Listen(":3000"))
}
