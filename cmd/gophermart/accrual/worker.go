package accrual

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/config/accrual"
)

type Worker struct {
	client *http.Client

	AccrualAddress string
	updateInterval time.Duration

	logger *zap.Logger
}

func New(cfg accrual.Config, logger *zap.Logger) *Worker {
	return &Worker{
		client:         &http.Client{}, //nolint:exhaustruct
		AccrualAddress: cfg.Address,

		updateInterval: cfg.UpdateInterval,

		logger: logger,
	}
}

func (w *Worker) Run(ctx context.Context, su OrdersSearcherUpdater) {
	ticker := time.NewTicker(w.updateInterval)

	select {
	case <-ticker.C:
		if err := w.processAccrual(su); err != nil { //nolint:contextcheck
			w.logger.Error("accrual processing failed", zap.Error(err))
		}

	case <-ctx.Done():
		return
	}
}
