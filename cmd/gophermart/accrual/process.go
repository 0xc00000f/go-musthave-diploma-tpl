package accrual

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/storage"
	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/structures/status"
)

type OrderUpdater interface {
	Update(ctx context.Context, data storage.OrderUpdateData) (*storage.OrderData, error)
}

type OrdersSearcher interface {
	Search(ctx context.Context, query storage.OrdersSearchQuery) (*storage.OrdersSearchResult, error)
}

type OrdersSearcherUpdater interface {
	OrdersSearcher
	OrderUpdater
}

func (w *Worker) processAccrual(su OrdersSearcherUpdater) error {
	orders, err := su.Search(context.Background(), storage.OrdersSearchQuery{
		Status: []status.OrderStatus{
			status.OrderStatusNew,
			status.OrderStatusProcessing,
		}})
	if err != nil {
		return fmt.Errorf("search orders failed: %w", err)
	}

	eg := &errgroup.Group{}
	for _, order := range orders.Orders {
		response, err := w.fetchInfo(order.OrderNumber)
		if err != nil {
			return fmt.Errorf("fetch info failed: %w", err)
		}

		eg.Go(func() error {
			if _, err := su.Update(context.Background(), storage.OrderUpdateData{
				OrderNumber: response.OrderNumber,
				Accrual:     response.Accrual,
				Status:      response.Status,
			}); err != nil {
				return fmt.Errorf("update order failed: %w", err)
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		w.logger.Error("error processing order", zap.Error(err))
		return err
	}

	return nil
}
