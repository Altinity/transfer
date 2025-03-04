package batchsplitter

import (
	"github.com/altinity/transfer/internal/logger"
	"github.com/altinity/transfer/library/go/core/metrics"
	"github.com/altinity/transfer/pkg/abstract"
	"github.com/altinity/transfer/pkg/abstract/coordinator"
	"github.com/altinity/transfer/pkg/abstract/model"
	"github.com/altinity/transfer/pkg/middlewares"
	"github.com/altinity/transfer/pkg/transformer"
	"github.com/altinity/transfer/pkg/util"
)

func PluggableBatchSplitterTransformer(transfer *model.Transfer, _ metrics.Registry, _ coordinator.Coordinator) func(abstract.Sinker) abstract.Sinker {
	if transfer.Transformation == nil || transfer.Transformation.Transformers == nil {
		return IdentityMiddleware
	}

	config := transferNeedDetector(transfer.Transformation.Transformers)
	if config == nil {
		return IdentityMiddleware
	}

	return func(s abstract.Sinker) abstract.Sinker {
		return newPluggableTransformer(s, *config)
	}
}

var IdentityMiddleware = func(s abstract.Sinker) abstract.Sinker { return s }

func transferNeedDetector(transformers *transformer.Transformers) *Config {
	for _, t := range transformers.Transformers {
		if v, ok := t[Type]; ok {
			var cfg Config
			if err := util.MapFromJSON(v, &cfg); err != nil {
				logger.Log.Errorf("unable to map %v to %v: %v", v, cfg, err)
				return nil
			}
			return &cfg
		}
	}
	return nil
}

type pluggableTransformer struct {
	sink   abstract.Sinker
	config Config
}

func newPluggableTransformer(s abstract.Sinker, cfg Config) abstract.Sinker {
	return &pluggableTransformer{s, cfg}
}

func (d *pluggableTransformer) Close() error {
	return d.sink.Close()
}

func (d *pluggableTransformer) Push(items []abstract.ChangeItem) error {
	for start := 0; start < len(items); start += d.config.MaxItemsPerBatch {
		end := start + d.config.MaxItemsPerBatch
		if end > len(items) {
			end = len(items)
		}
		err := d.sink.Push(items[start:end])
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	middlewares.PlugTransformer(PluggableBatchSplitterTransformer)
}
