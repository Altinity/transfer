package tasks

import (
	"context"

	"github.com/altinity/transfer/internal/logger"
	"github.com/altinity/transfer/library/go/core/metrics"
	"github.com/altinity/transfer/pkg/abstract/coordinator"
	"github.com/altinity/transfer/pkg/abstract/model"
	"github.com/altinity/transfer/pkg/providers"
)

func Deactivate(ctx context.Context, cp coordinator.Coordinator, transfer model.Transfer, task model.TransferOperation, registry metrics.Registry) error {
	deactivator, ok := providers.Source[providers.Deactivator](logger.Log, registry, cp, &transfer)
	if ok {
		return deactivator.Deactivate(ctx, &task)
	}
	return nil
}
