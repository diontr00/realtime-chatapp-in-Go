package route

import (
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"{{{template}}}/api/controller"
	"{{{template}}}/api/middleware"
	"{{{template}}}/config"
	"{{{template}}}/translator"
)

type RouteConfig struct {
	Env        *config.Env
	Timeout    time.Duration
	Fiber      *fiber.App
	Validator  *validator.Validate
	Translator *translator.UTtrans
}

func Setup(config *RouteConfig) {

	main_controller := &controller.Maincontroller{
		Validator:  config.Validator,
		Translator: config.Translator,
	}
	config.Fiber.Use(middleware.NewLocaleMiddleWare)
	config.Fiber.Get("/", main_controller.Hello)
}
