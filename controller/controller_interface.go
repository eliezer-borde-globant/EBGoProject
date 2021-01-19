package controller

import "github.com/gofiber/fiber/v2"

var (
	ControllerObject controllerInterface = controllerImplementation{}
)

type controllerInterface interface {
	UpdateSecretfile(c *fiber.Ctx) error
	CreateSecretfile(c *fiber.Ctx) error
}

type controllerImplementation struct { }