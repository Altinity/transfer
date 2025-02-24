package targets

import (
	"github.com/altinity/transfer/library/go/core/metrics"
	"github.com/altinity/transfer/library/go/core/xerrors"
	"github.com/altinity/transfer/pkg/abstract"
	"github.com/altinity/transfer/pkg/abstract/coordinator"
	"github.com/altinity/transfer/pkg/abstract/model"
	"github.com/altinity/transfer/pkg/base"
	"github.com/altinity/transfer/pkg/providers"
	"go.ytsaurus.tech/library/go/core/log"
)

var UnknownTargetError = xerrors.New("unknown event target for destination, try legacy sinker instead")

func NewTarget(transfer *model.Transfer, lgr log.Logger, mtrcs metrics.Registry, cp coordinator.Coordinator, opts ...abstract.SinkOption) (t base.EventTarget, err error) {
	if factory, ok := providers.Destination[providers.Abstract2Sinker](lgr, mtrcs, cp, transfer); ok {
		return factory.Target(opts...)
	}
	return nil, UnknownTargetError
}
