package typefitting

import (
	"github.com/altinity/transfer/pkg/abstract/model"
	chrecipe "github.com/altinity/transfer/pkg/providers/clickhouse/recipe"
)

var (
	//nolint:exhaustivestruct
	source = model.MockSource{}
	target = *chrecipe.MustTarget(chrecipe.WithDatabase("test"), chrecipe.WithInitFile("init.sql"))
)

func init() {
	source.WithDefaults()
	target.WithDefaults()
}
