package coingeckoclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/btc-price/pkg/btcpricelb"
)

type Client struct {
	client   *http.Client
	ratePath string
}

func NewClient(path string, cl *http.Client) *Client {
	return &Client{
		client:   cl,
		ratePath: path,
	}
}

func (c *Client) GetRate(ctx context.Context, marketCurr string, baseCurr string) (btcpricelb.CoingeckoRate, error) {
	q := url.Values{}
	q.Set("ids", marketCurr)
	q.Set("vs_currencies", baseCurr)
	path := fmt.Sprintf("%s?%s", c.ratePath, q.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
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
