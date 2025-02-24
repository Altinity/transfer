package events

import (
	"github.com/altinity/transfer/pkg/base"
	"github.com/altinity/transfer/pkg/providers/yt/tablemeta"
)

type tableEvent struct {
	path *tablemeta.YtTableMeta
}

type TableEvent interface {
	base.Event
	Table() *tablemeta.YtTableMeta
}

func (t *tableEvent) Table() *tablemeta.YtTableMeta {
	return t.path
}

func newTableEvent(path *tablemeta.YtTableMeta) TableEvent {
	return &tableEvent{path}
}
