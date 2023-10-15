package main

import (
	"context"
	"sync"

	"realtime-chat/api/route"
	"realtime-chat/config"
)

var (
	app_    *config.Applications
	appOnce sync.Once
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	appOnce.Do(func() {
		app_ = config.NewApp(ctx)
	})

	route.Setup(ctx, &route.RouteConfig{
		Env:        app_.Env,
		Timeout:    app_.Env.App.Timeout,
		Validator:  app_.Validator,
		Translator: app_.Translator,
		Fiber:      app_.Fiber,
	})

	app_.Fiber.Static("/", "./frontend/dist")

	app_.Start()
}
