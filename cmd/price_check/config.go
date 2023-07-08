package main

import "time"

type Config struct {
	Port          string        `env:"PORT"           envDefault:":8080"`
	ServerTimeout time.Duration `env:"SERVER_TIMEOUT" envDefault:"10s"`
	Market        Market
	Coingecko     Coingecko
	Binance       Binance
	EmailStorage  FileStorage
	EmailText     EmailText
}

type Market struct {
	BaseCurrency  string `env:"BASE_CURRENCY"  envDefault:"BTC"`
	QuoteCurrency string `env:"QUOTE_CURRENCY" envDefault:"UAH"`
}

type Coingecko struct {
	RatePath string `env:"COINGECKO_RATE_PATH" envDefault:"https://api.coingecko.com/api/v3/simple/price"`
}

type FileStorage struct {
	Path string `env:"FILE_PATH" envDefault:"./emails.txt"`
}

type Binance struct {
	RatePath string `env:"BINANCE_RATE_PATH" envDefault:"https://api.binance.com/api/v3/avgPrice"`
}

type EmailText struct {
	Text string `env:"TEXT" envDefault:"Current rate is"`
}
