package emailsender

import (
	"context"

	"github.com/btc-price/pkg/btcpricelb"
)

type Sender struct {
}

func NewSender() *Sender {
	return &Sender{}
}

func (s *Sender) SendEmails(_ context.Context, _ []btcpricelb.Email) error {
	// interaction with third-party email sending service
	// for _, _ = range emailsList {
	// }
	return nil
}
