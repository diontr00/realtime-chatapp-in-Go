package middleware

import (
	"github.com/gofiber/fiber/v2"
	"{{{template}}}/translator"
)

func NewCatchAllMiddleWare(trans translator.UTtrans) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Error": trans.TranslateMessage(c, "unknownroute", translator.TranslateParam{"Method": c.Method(), "Route": c.Route().Name}, nil),
		})
	}

}
