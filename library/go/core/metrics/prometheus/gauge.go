package prometheus

import (
	"github.com/altinity/transfer/library/go/core/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

var _ metrics.Gauge = (*Gauge)(nil)

// Gauge tracks single float64 value.
type Gauge struct {
	gg prometheus.Gauge
}

func (g Gauge) Set(value float64) {
	g.gg.Set(value)
}

func (g Gauge) Add(value float64) {
	g.gg.Add(value)
}

var _ metrics.FuncGauge = (*FuncGauge)(nil)

type FuncGauge struct {
	ff       prometheus.GaugeFunc
	function func() float64
}

func (g FuncGauge) Function() func() float64 {
	return g.function
}
