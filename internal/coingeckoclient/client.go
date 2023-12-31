package coingeckoclient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/btc-price/pkg/coinconverter"
	"io"
	"net/http"
	"net/url"

	"github.com/btc-price/pkg/btcpricelb"
)

type Client struct {
	client    *http.Client
	ratePath  string
	converter *coinconverter.Converter
}

func NewClient(cl *http.Client, path string, converter *coinconverter.Converter) *Client {
	return &Client{
		client:    cl,
		ratePath:  path,
		converter: converter,
	}
}

func (c *Client) GetRateRequest(ctx context.Context, bCurr string, qCurr string) (btcpricelb.Rate, error) {
	q := url.Values{}
	q.Set("ids", c.converter.ConvertCoingeckoBase(bCurr))
	q.Set("vs_currencies", c.converter.ConvertCoingeckoQuote(qCurr))
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

	return btcpricelb.Rate(answer.Bitcoin.Uah), nil
}
