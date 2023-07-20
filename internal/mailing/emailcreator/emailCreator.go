package emailcreator

import (
	"context"
	"fmt"
	"github.com/btc-price/pkg/btcpricelb"
)

type EmailCreator struct {
	text string
}

func NewEmailCreator(text string) *EmailCreator {
	return &EmailCreator{text: text}
}

func (e *EmailCreator) GenerateEmail(_ context.Context, rate btcpricelb.Rate) string {
	return fmt.Sprintf("%s %f", e.text, rate)
}
