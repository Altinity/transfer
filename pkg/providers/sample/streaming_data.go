package sample

import (
	"github.com/altinity/transfer/pkg/abstract"
)

type StreamingData interface {
	TableName() abstract.TableID
	ToChangeItem(offset int64) abstract.ChangeItem
}
