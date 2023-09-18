package route

import (
	"time"

	"github.com/diontr00/distributedkv/api/controller"
	"github.com/diontr00/distributedkv/api/middleware"
	"github.com/diontr00/distributedkv/config"
	"github.com/diontr00/distributedkv/translator"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
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
