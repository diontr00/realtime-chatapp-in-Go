package route

import (
	"context"
	"time"

	"realtime-chat/api/controller"
	"realtime-chat/api/middleware"
	"realtime-chat/config"
	"realtime-chat/translator"

	"github.com/go-playground/validator"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	Env        *config.Env
	Timeout    time.Duration
	Fiber      *fiber.App
	Validator  *validator.Validate
	Translator *translator.UTtrans
}

func Setup(ctx context.Context, config *RouteConfig) {

	// main_controller := controller.NewMainController(config.Validator, config.Translator)

	ws_controller := controller.NewSocketController(ctx, config.Translator, config.Validator, config.Env.Socket)

	// Middleware
	config.Fiber.Use(middleware.NewLocaleMiddleWare)
	config.Fiber.Use(middleware.NewRecoverMiddleWare(config.Translator))
	config.Fiber.Use("/ws", middleware.Newwsmiddleware)

	// Route

	config.Fiber.Static("/", "./frontend/dist")

	config.Fiber.Get("/ws", websocket.New(ws_controller.Serve, websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		Origins:         []string{"https://localhost:8080"},
	}))

	config.Fiber.Post("/login", ws_controller.LoginHandler)

	config.Fiber.All("*", middleware.NewCatchAllMiddleWare(config.Translator))

}
