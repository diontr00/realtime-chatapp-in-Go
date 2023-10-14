package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"realtime-chat/translator"
)

func NewRecoverMiddleWare(trans *translator.UTtrans) func(c *fiber.Ctx) error {
	return recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			log.Printf("[PANIC] : recovered ! 😥 -- %v", e)
			log.Println(c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"Error": trans.TranslateMessage(c.Locals("locale").(string), "internal", nil, nil),
			}))

		},
	})
}
