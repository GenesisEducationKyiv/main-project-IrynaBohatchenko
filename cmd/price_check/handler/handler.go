package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/btc-price/internal/storageerrors"

	"github.com/btc-price/pkg/btcpricelb"
	"go.uber.org/zap"
)

type (
	BtcPrice struct {
		srv    BtcPriceService
		logger *zap.Logger
	}

	BtcPriceService interface {
		HandleRate(ctx context.Context, marketCurr string, baseCurr string) (btcpricelb.RateResponse, error)
		HandleSubscribe(ctx context.Context, email string) error
		HandleSendEmails(ctx context.Context) error
	}
)

func NewBtcPrice(sr BtcPriceService, logger *zap.Logger) *BtcPrice {
	return &BtcPrice{
		srv:    sr,
		logger: logger,
	}
}

func (b *BtcPrice) handleRate(writer http.ResponseWriter, request *http.Request) {
	logger := b.logger.Named("rate handler")

	resp, err := b.srv.HandleRate(request.Context(), "bitcoin", "uah")
	if err != nil {
		logger.Error("error getting rate", zap.Error(err))
		b.write(writer, http.StatusBadRequest, "error getting rate")

		return
	}

	b.write(writer, http.StatusOK, resp.Rate)
}

func (b *BtcPrice) handleSubscribe(writer http.ResponseWriter, request *http.Request) {
	logger := b.logger.Named("subscription handler")

	email := request.FormValue(btcpricelb.EmailForm)

	if err := b.srv.HandleSubscribe(request.Context(), email); err != nil {
		logger.Error("error subscribing email", zap.Error(err))

		status := http.StatusInternalServerError
		errText := btcpricelb.RespTextSubscrErr
		if errors.Is(err, storageerrors.ErrEmailExists) {
			status = http.StatusConflict
			errText = btcpricelb.RespTextEmailExists
		}

		b.write(writer, status, errText)

		return
	}

	b.write(writer, http.StatusOK, "E-mail додано")
}

func (b *BtcPrice) handleSendEmails(writer http.ResponseWriter, request *http.Request) {
	logger := b.logger.Named("send emails handler")

	if err := b.srv.HandleSendEmails(request.Context()); err != nil {
		logger.Error("error sending emails", zap.Error(err))

		return
	}

	b.write(writer, http.StatusOK, "E-mailʼи відправлено")
}

func (b *BtcPrice) write(writer http.ResponseWriter, statusCode int, data interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	if bt, ok := data.([]byte); ok {
		if _, err := writer.Write(bt); err != nil {
			b.logger.Error(", write response", zap.Error(err), zap.Any("response", data))

			return
		}
	}

	b.logger.Info("writer.go data=", zap.Any("data", data))
	if err := json.NewEncoder(writer).Encode(data); err != nil {
		b.logger.Error("json encoder, write response", zap.Error(err), zap.Any("response", data))
	}
}
