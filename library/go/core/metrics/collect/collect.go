package collect

import (
	"context"

	"github.com/altinity/transfer/library/go/core/metrics"
)

type Func func(ctx context.Context, r metrics.Registry, c metrics.CollectPolicy)
