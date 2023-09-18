package controller

import (
	"github.com/diontr00/distributedkv/translator"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
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
