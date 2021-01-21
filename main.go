package main

import (
	"fmt"

	"github.com/eliezer-borde-globant/EBGoProject/controller"
	"github.com/eliezer-borde-globant/EBGoProject/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
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
	data := new(controller.CreateParams)
	if err := c.BodyParser(data); err != nil {
		utils.ZeroLogger.Error().Msgf("Check input data , error: %v", err)
		return c.Status(400).SendString(fmt.Sprintf("contents not parsed correctly: %v", err))
	}
	statusCode, msg := controller.ControllerObject.CreateSecretFile(data)
	return c.Status(statusCode).SendString(fmt.Sprintf(msg))

}

func updateSecretFile(c *fiber.Ctx) error {
	data := new(controller.UpdateParams)
	if err := c.BodyParser(data); err != nil {
		utils.ZeroLogger.Error().Msgf("Check input data , error: %v", err)
		return c.Status(400).SendString(fmt.Sprintf("contents not parsed correctly: %v", err))
	}
	statusCode, msg := controller.ControllerObject.UpdateSecretFile(data)
	return c.Status(statusCode).SendString(fmt.Sprintf(msg))

}
