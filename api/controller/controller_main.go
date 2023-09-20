package controller

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"{{{template}}}/translator"
)

type Maincontroller struct {
	Validator  *validator.Validate
	Translator *translator.UTtrans
}

func (m *Maincontroller) Hello(c *fiber.Ctx) error {

	message := m.Translator.TranslateMessage(c, "hello", nil, nil)

	return c.JSON(fiber.Map{
		"message": message,
	})

}
