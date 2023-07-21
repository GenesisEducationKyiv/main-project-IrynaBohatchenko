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
		rateSrv   RateService
		subscrSrv SubscriptionService
		sendSrv   SendService
		logger    *zap.Logger
	}

	RateService interface {
		HandleRate(ctx context.Context, bCurr string, qCurr string) (btcpricelb.RateResponse, error)
	}

	SubscriptionService interface {
		HandleSubscribe(ctx context.Context, email string) error
	}

	SendService interface {
		HandleSendEmails(ctx context.Context, bCurr, qCurr string) error
	}
)

func NewBtcPrice(
	rateSrv RateService,
	sbscrSrv SubscriptionService,
	sndSrv SendService,
	logger *zap.Logger) *BtcPrice {
	return &BtcPrice{
		rateSrv:   rateSrv,
		subscrSrv: sbscrSrv,
		sendSrv:   sndSrv,
		logger:    logger,
	}
}

func (b *BtcPrice) handleRate(writer http.ResponseWriter, request *http.Request) {
	logger := b.logger.Named("rate handler")

	resp, err := b.rateSrv.HandleRate(request.Context(), "bitcoin", "uah")
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

	if err := b.subscrSrv.HandleSubscribe(request.Context(), email); err != nil {
		logger.Error("error subscribing email", zap.Error(err))

		status := http.StatusInternalServerError
		errText := btcpricelb.RespTextSubscrErr
		switch {
		case errors.Is(err, storageerrors.ErrEmailExists):
			status = http.StatusConflict
			errText = btcpricelb.RespTextEmailExists
		case errors.Is(err, storageerrors.ErrInvalidEmail):
			status = http.StatusBadRequest
			errText = btcpricelb.RespTextInvalidEmail
		}

		b.write(writer, status, errText)

		return
	}

	b.write(writer, http.StatusOK, "E-mail додано")
}

func (b *BtcPrice) handleSendEmails(writer http.ResponseWriter, request *http.Request) {
	logger := b.logger.Named("send emails handler")

	if err := b.sendSrv.HandleSendEmails(request.Context(), "bitcoin", "uah"); err != nil {
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
