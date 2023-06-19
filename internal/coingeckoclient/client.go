package coingeckoclient

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/btc-price/pkg/btcpricelb"
)

type Client struct {
	client   *http.Client
	ratePath string
}

func NewClient(path string) *Client {
	return &Client{
		client:   &http.Client{},
		ratePath: path,
	}
}

func (c *Client) GetRate(ctx context.Context) (btcpricelb.CoingeckoRate, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.ratePath, nil)
	if err != nil {
		return 0, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	answerByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var answer btcpricelb.CoingeckoResponse
	if err := json.Unmarshal(answerByte, &answer); err != nil {
		return 0, err
	}

	return btcpricelb.CoingeckoRate(answer.Bitcoin.Uah), nil
}
