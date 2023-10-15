package config

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"realtime-chat/translator"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type Applications struct {
	Env        *Env
	Fiber      *fiber.App
	Translator *translator.UTtrans
	Validator  *validator.Validate
}

func (a *Applications) Start() {
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)

	go func() {
		<-terminate
		fmt.Printf("Gratefully Shutdown %s , Doing Cleanup Task...😷 \n", a.Env.Fiber.AppName)
		log.Fatal(a.Fiber.Shutdown())
	}()

	log.Fatal(a.Fiber.ListenTLS(a.Env.Fiber.ListenPort, "./cert/localhost.cert", "./cert/localhost.key"))
}

func NewApp(ctx context.Context) *Applications {

	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	env := newEnv(ctx)
	server := newFiber(env.Fiber)
	trans := newTranslator()
	valtor := newValidator()

	return &Applications{
		Env:        env,
		Fiber:      server,
		Translator: trans,
		Validator:  valtor,
	}

}

func newFiber(config FiberEnv) *fiber.App {
	server := fiber.New(fiber.Config{AppName: config.AppName})
	return server

}

//go:embed trans_file/*.toml
var trans_folder embed.FS

func newTranslator() *translator.UTtrans {
	trans, err := translator.NewUtTrans(trans_folder, "trans_file")
	if err != nil {
		log.Fatalf("[Error] Reading Translation File %v \n", err)
	}
	return trans
}

func newValidator() *validator.Validate {
	return validator.New()
}
