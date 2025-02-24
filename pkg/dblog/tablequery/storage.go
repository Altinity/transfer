package tablequery

import (
	"context"

	"github.com/altinity/transfer/pkg/abstract"
)

// StorageTableQueryable is storage with table query loading
type StorageTableQueryable interface {
	abstract.SampleableStorage

	LoadQueryTable(ctx context.Context, table TableQuery, pusher abstract.Pusher) error
}
