package opensearch

import (
	"context"

	"github.com/altinity/transfer/pkg/abstract"
)

func (s *Storage) ShardTable(ctx context.Context, table abstract.TableDescription) ([]abstract.TableDescription, error) {
	return s.elasticShardingStorage.ShardTable(ctx, table)
}
