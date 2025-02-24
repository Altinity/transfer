package helpers

import (
	"context"
	"testing"

	"github.com/altinity/transfer/internal/logger"
	"github.com/altinity/transfer/library/go/core/xerrors"
	"github.com/altinity/transfer/pkg/abstract"
	"github.com/altinity/transfer/pkg/abstract/coordinator"
	"github.com/altinity/transfer/pkg/abstract/model"
	"github.com/altinity/transfer/pkg/providers/postgres"
	"github.com/altinity/transfer/pkg/runtime/local"
	"github.com/altinity/transfer/pkg/worker/tasks"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

type Worker struct {
	worker *local.LocalWorker
	cp     coordinator.Coordinator
}

// controlplane that catches replication failure
type fakeCpErrRepl struct {
	coordinator.Coordinator
	onErrorCallback []func(err error)
}

func (f *fakeCpErrRepl) FailReplication(transferID string, err error) error {
	for _, cb := range f.onErrorCallback {
		cb(err)
	}
	return nil
}

func (q *Worker) Close(t *testing.T) {
	if q.worker != nil {
		err := q.worker.Stop()
		if xerrors.Is(err, context.Canceled) {
			return
		}
		require.NoError(t, err)
	}
}

// Restart replication worker with updated transfer
func (q *Worker) Restart(t *testing.T, transfer *model.Transfer) {
	q.Close(t)
	q.initLocalWorker(transfer)
	q.worker.Start()
}

func (q *Worker) initLocalWorker(transfer *model.Transfer) {
	q.worker = local.NewLocalWorker(q.cp, transfer, EmptyRegistry(), logger.LoggerWithLevel(zapcore.DebugLevel))
}

func Activate(t *testing.T, transfer *model.Transfer, onErrorCallback ...func(err error)) *Worker {
	if len(onErrorCallback) == 0 {
		// append default callback checker: no error!
		onErrorCallback = append(onErrorCallback, func(err error) {
			require.NoError(t, err)
		})
	}
	result, err := ActivateErr(transfer, onErrorCallback...)
	require.NoError(t, err)

	return result
}

func ActivateErr(transfer *model.Transfer, onErrorCallback ...func(err error)) (*Worker, error) {
	cp := &fakeCpErrRepl{Coordinator: coordinator.NewStatefulFakeClient(), onErrorCallback: onErrorCallback}
	return ActivateWithCP(transfer, cp)
}

func ActivateWithCP(transfer *model.Transfer, cp coordinator.Coordinator) (*Worker, error) {
	result := &Worker{
		worker: nil,
		cp:     cp,
	}
	err := tasks.ActivateDelivery(context.Background(), nil, result.cp, *transfer, EmptyRegistry())
	if err != nil {
		return nil, err
	}

	if transfer.Type == abstract.TransferTypeSnapshotAndIncrement || transfer.Type == abstract.TransferTypeIncrementOnly {
		if pgDst, ok := transfer.Dst.(*postgres.PgDestination); ok {
			pgDst.CopyUpload = false
		}
		result.initLocalWorker(transfer)
		result.worker.Start()
	}

	return result, nil
}
