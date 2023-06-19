package main

import "time"

type Config struct {
	Port          string        `env:"PORT"           envDefault:":8080"`
	ServerTimeout time.Duration `env:"SERVER_TIMEOUT" envDefault:"10s"`
	ClientTimeout time.Duration `env:"CLIENT_TIMEOUT" envDefault:"30s"`
	Coingecko     Coingecko
	EmailStorage  FileStorage
}

type Coingecko struct {
	RatePath string `env:"RATE_PATH" envDefault:"https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=uah"` //nolint:lll
}

type FileStorage struct {
	Path string `env:"PATH" envDefault:"./emails.txt"`
}
