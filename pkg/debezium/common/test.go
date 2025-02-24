package common

import "github.com/altinity/transfer/pkg/abstract"

type ChangeItemCanon struct {
	ChangeItem     *abstract.ChangeItem
	DebeziumEvents []KeyValue
}
