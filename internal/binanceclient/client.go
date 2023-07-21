package binanceclient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/btc-price/pkg/btcpricelb"
	"github.com/btc-price/pkg/coinconverter"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Client struct {
	client    *http.Client
	ratePath  string
	converter *coinconverter.Converter
}

func NewClient(client *http.Client, ratePath string, converter *coinconverter.Converter) *Client {
	return &Client{
		client:    client,
		ratePath:  ratePath,
		converter: converter,
	}
}

func (c *Client) GetRateRequest(ctx context.Context, bCurr string, qCurr string) (btcpricelb.Rate, error) {
	symbol := strings.ToUpper(fmt.Sprint(c.converter.ConvertBinanceCoins(bCurr), c.converter.ConvertBinanceCoins(qCurr)))
	q := url.Values{}
	q.Set("symbol", symbol)
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

	var answer btcpricelb.BinanceResponse
	if err := json.Unmarshal(answerByte, &answer); err != nil {
		return 0, err
	}

	rate, err := strconv.ParseFloat(answer.Price, 64)
	if err != nil {
		return 0, err
	}

	return btcpricelb.Rate(rate), nil
}
