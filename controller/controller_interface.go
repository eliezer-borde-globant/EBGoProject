package controller

import "github.com/gofiber/fiber/v2"

var (
	ControllerObject controllerInterface = controllerImplementation{}
)

type controllerInterface interface {
	UpdateSecretFile(c *fiber.Ctx) error
	CreateSecretFile(createInterface contextCreateInterface) error
}

type contextCreateInterface interface {
	BodyParser(data *createParams) error
	Status(code int) *fiber.Ctx
}

type controllerImplementation struct { }