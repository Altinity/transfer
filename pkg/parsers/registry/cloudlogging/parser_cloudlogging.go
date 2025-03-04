package cloudlogging

import (
	"github.com/altinity/transfer/pkg/parsers"
	cloudloggingengine "github.com/altinity/transfer/pkg/parsers/registry/cloudlogging/engine"
	"github.com/altinity/transfer/pkg/stats"
	"go.ytsaurus.tech/library/go/core/log"
)

func NewParserCloudLogging(inWrapped interface{}, sniff bool, logger log.Logger, registry *stats.SourceStats) (parsers.Parser, error) {
	return cloudloggingengine.NewCloudLoggingImpl(sniff, logger, registry), nil
}

func init() {
	parsers.Register(
		NewParserCloudLogging,
		[]parsers.AbstractParserConfig{new(ParserConfigCloudLoggingCommon)},
	)
}
