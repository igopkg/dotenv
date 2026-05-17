package main

import (
	"fmt"
	"log"
	"time"

	"github.com/igopkg/dotenv"
)

type Config struct {
	Server ServerConfig
	DB     DBConfig
}

type ServerConfig struct {
	Addr    string        `env:"SERVER_ADDR"    req:"true"`
	Timeout time.Duration `env:"SERVER_TIMEOUT"`
	Debug   bool          `env:"DEBUG"`
}

type DBConfig struct {
	Host     string `env:"DB_HOST" req:"true"`
	Port     string `env:"DB_PORT" req:"true"`
	Name     string `env:"DB_NAME" req:"true"`
	Password string `env:"DB_PASSWORD"`
}

func main() {
	// Step 1: load .env file into the process environment.
	if err := dotenv.Load("./example/env/.env"); err != nil {
		log.Fatalf("load .env: %v", err)
	}

	// Step 2: unmarshal environment variables into the config struct.
	var cfg Config
	if err := dotenv.Unmarshal(&cfg); err != nil {
		log.Fatalf("unmarshal config: %v", err)
	}

	fmt.Printf("server addr:    %s\n", cfg.Server.Addr)
	fmt.Printf("server timeout: %s\n", cfg.Server.Timeout)
	fmt.Printf("debug:          %v\n", cfg.Server.Debug)
	fmt.Printf("db host:        %s\n", cfg.DB.Host)
	fmt.Printf("db port:        %s\n", cfg.DB.Port)
	fmt.Printf("db name:        %s\n", cfg.DB.Name)
}
