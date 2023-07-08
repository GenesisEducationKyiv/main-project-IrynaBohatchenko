package btcpricelb

type RateResponse struct {
	Rate float64 `json:"rate"`
}

type Rate float64

type CoingeckoResponse struct {
	Bitcoin struct {
		Uah float64 `json:"uah"`
	} `json:"bitcoin"`
}

type BinanceResponse struct {
	Price string `json:"price"`
}
