package confluentschemaregistry

import (
	"github.com/altinity/transfer/pkg/parsers"
	conflueentschemaregistryengine "github.com/altinity/transfer/pkg/parsers/registry/confluentschemaregistry/engine"
	"github.com/altinity/transfer/pkg/stats"
	"go.ytsaurus.tech/library/go/core/log"
)

func NewParserConfluentSchemaRegistry(inWrapped interface{}, _ bool, logger log.Logger, _ *stats.SourceStats) (parsers.Parser, error) {
	switch in := inWrapped.(type) {
	case *ParserConfigConfluentSchemaRegistryCommon:
		return conflueentschemaregistryengine.NewConfluentSchemaRegistryImpl(in.SchemaRegistryURL, in.TLSFile, in.Username, in.Password, false, logger), nil
	case *ParserConfigConfluentSchemaRegistryLb:
		return conflueentschemaregistryengine.NewConfluentSchemaRegistryImpl(in.SchemaRegistryURL, in.TLSFile, in.Username, in.Password, false, logger), nil
	}
	return nil, nil
}

func init() {
	parsers.Register(
		NewParserConfluentSchemaRegistry,
		[]parsers.AbstractParserConfig{new(ParserConfigConfluentSchemaRegistryCommon), new(ParserConfigConfluentSchemaRegistryLb)},
	)
}
