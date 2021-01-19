package main

import (
	. "github.com/eliezer-borde-globant/EBGoProject/controller"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Static("/", "./public")
	app.Post("/api/detectsecrets/update", ControllerObject.UpdateSecretFile)
	app.Post("/api/detectsecrets/create", ControllerObject.CreateSecretFile)
	app.Listen(":3000")
}