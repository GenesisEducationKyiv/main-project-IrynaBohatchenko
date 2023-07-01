package btcpricelb

type RateResponse struct {
	Rate float64 `json:"rate"`
}

type CoingeckoRate float64

type CoingeckoResponse struct {
	Bitcoin struct {
		Uah float64 `json:"uah"`
	} `json:"bitcoin"`
}

type Email string
