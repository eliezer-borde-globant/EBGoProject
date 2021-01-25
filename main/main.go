package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

var (
	controllerObjectCreateSecretFile = ControllerObject.CreateSecretFile
	controllerObjectUpdateSecretFile = ControllerObject.UpdateSecretFile
)

func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Static("/", "./public")
	app.Post("/api/detectsecrets/update", updateSecretFile)
	app.Post("/api/detectsecrets/create", createSecretFile)
	app.Listen(":3000")
}

func createSecretFile(c *fiber.Ctx) error {
	data := new(CreateParams)
	if err := c.BodyParser(data); err != nil {
		zeroLogger.Error().Msgf("Check input data , error: %v", err)
		return c.Status(400).SendString(fmt.Sprintf("contents not parsed correctly: %v", err))
	}
	statusCode, msg := controllerObjectCreateSecretFile(data)
	return c.Status(statusCode).SendString(fmt.Sprintf(msg))

}

func updateSecretFile(c *fiber.Ctx) error {
	data := new(UpdateParams)
	if err := c.BodyParser(data); err != nil {
		zeroLogger.Error().Msgf("Check input data , error: %v", err)
		return c.Status(400).SendString(fmt.Sprintf("contents not parsed correctly: %v", err))
	}

	statusCode, msg := controllerObjectUpdateSecretFile(data)
	return c.Status(statusCode).SendString(fmt.Sprintf(msg))

}
