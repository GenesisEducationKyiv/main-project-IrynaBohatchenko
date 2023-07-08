// Package coinconverter TODO refactor to stop break SRP and Open-closed
package coinconverter

var coingeckoBaseCoins = map[string]string{
	"BTC": "bitcoin",
}

var coingeckoQuoteCoins = map[string]string{
	"UAH": "uah",
}

var binanceCoins = map[string]string{
	"BTC": "BTC",
	"UAH": "UAH",
}

type Converter struct {
	coingeckoBase  map[string]string
	coingeckoQuote map[string]string
	binanceCoins   map[string]string
}

func NewConverter() *Converter {
	return &Converter{
		coingeckoBase:  coingeckoBaseCoins,
		coingeckoQuote: coingeckoQuoteCoins,
		binanceCoins:   binanceCoins,
	}
}

func (c *Converter) ConvertCoingeckoBase(coin string) string {
	return c.coingeckoBase[coin]
}

func (c *Converter) ConvertCoingeckoQuote(coin string) string {
	return c.coingeckoQuote[coin]
}

func (c *Converter) ConvertBinanceCoins(coin string) string {
	return c.binanceCoins[coin]
}
