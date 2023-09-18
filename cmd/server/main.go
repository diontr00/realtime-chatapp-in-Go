package main

import (
	"context"
	"sync"
	"time"

	"github.com/diontr00/distributedkv/api/route"
	"github.com/diontr00/distributedkv/config"
)

var (
	app_    *config.Applications
	appOnce sync.Once
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	appOnce.Do(func() {
		app_ = config.NewApp(ctx)
	})

	route.Setup(&route.RouteConfig{
		Env:        app_.Env,
		Timeout:    app_.Env.App.Timeout,
		Validator:  app_.Validator,
		Translator: app_.Translator,
		Fiber:      app_.Fiber,
	})
	app_.Start()
}
