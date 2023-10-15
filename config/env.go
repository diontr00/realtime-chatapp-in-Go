package config

import (
	"context"
	"embed"
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

//go:embed dotenv/*
var dotenvFS embed.FS

type (
	Env struct {
		App    AppEnv
		Socket SocketEnv
		Fiber  FiberEnv
	}

	AppEnv struct {
		Timeout    time.Duration `env:"TIMEOUT,default=2s"`
		Production bool          `env:"PRODUCTION,default=false"`
	}

	SocketEnv struct {
		PongWait     time.Duration `env:"PONGWAIT,default=10s"`
		OtpRetention time.Duration `env:"OTPRETENTION,default=5s"`
		PingInterval time.Duration
	}

	FiberEnv struct {
		AppName    string `env:"APP_NAME,default=api-server"`
		ListenPort string `env:"LISTEN_PORT,default=:8080"`
	}
)

func newEnv(ctx context.Context) *Env {
	env := loadEnvFile(ctx)
	if !env.App.Production {
		fmt.Println("Running App in Development Env 🔥")
	}

	// Interval to send ping
	env.Socket.PingInterval = (env.Socket.PongWait * 9) / 10

	return env

}

// Load and process *.env into struct
func loadEnvFile(ctx context.Context) *Env {

	env_files := readEnvFiles("dotenv")

	err := godotenv.Load(env_files...)
	if err != nil {
		log.Fatalf("[Error] - Load env file, %v", err.Error())
	}

	env := &Env{}

	err = envconfig.Process(ctx, env)
	if err != nil {
		log.Fatalf("[Error] - serialize env file, %v", err.Error())
	}

	return env
}

// read nested embedded dotenv file  , create tempFile for goDotEnv to load
func readEnvFiles(path string) []string {
	env_files := []string{}
	var walkDir func(string)

	tmp := os.TempDir()
	walkDir = func(path string) {
		files, err := dotenvFS.ReadDir(path)
		if err != nil {
			log.Fatalf("[Error] - Read nested dotenv Dir : %v", err)
		}

		for _, file := range files {
			filePath := filepath.Join(path, file.Name())
			if file.IsDir() {
				walkDir(filePath)
			} else {
				fileName := filepath.Join(tmp, file.Name())
				data, err := dotenvFS.ReadFile(filePath)
				if err != nil {
					log.Fatalf("[Error] - Read embedded Env file , %v", err.Error())
				}
				err = os.WriteFile(fileName, data, 0600)

				if err != nil {
					log.Fatalf("[Error] - Write Temp Env file , %v", err.Error())
				}

				env_files = append(env_files, fileName)

			}

		}
	}

	walkDir(path)
	return env_files

}
