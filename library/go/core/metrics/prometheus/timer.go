package prometheus

import (
	"time"

	"github.com/altinity/transfer/library/go/core/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

var _ metrics.Timer = (*Timer)(nil)

// Timer measures gauge duration.
type Timer struct {
	gg prometheus.Gauge
}

func (t Timer) RecordDuration(value time.Duration) {
	t.gg.Set(value.Seconds())
}
