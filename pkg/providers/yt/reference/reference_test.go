package reference

import (
	"os"
	"testing"
	"time"

	"github.com/altinity/transfer/internal/logger"
	"github.com/altinity/transfer/library/go/core/metrics/solomon"
	"github.com/altinity/transfer/pkg/abstract"
	"github.com/altinity/transfer/pkg/abstract/coordinator"
	"github.com/altinity/transfer/pkg/providers/yt"
	"github.com/altinity/transfer/pkg/providers/yt/sink"
	"github.com/altinity/transfer/tests/canon/reference"
	"github.com/stretchr/testify/require"
)

func TestPushReferenceTable(t *testing.T) {
	Destination := &yt.YtDestination{
		Path:                "//home/cdc/tests/reference",
		Cluster:             os.Getenv("YT_PROXY"),
		CellBundle:          "default",
		PrimaryMedium:       "default",
		Static:              true,
		DisableDatetimeHack: true,
	}
	cfg := yt.NewYtDestinationV1(*Destination)
	cfg.WithDefaults()
	t.Run("static", func(t *testing.T) {
		sinker, err := sink.NewSinker(cfg, "", 0, logger.Log, solomon.NewRegistry(solomon.NewRegistryOpts()), coordinator.NewFakeClient(), nil)
		require.NoError(t, err)

		require.NoError(t, sinker.Push([]abstract.ChangeItem{
			{Kind: abstract.InitTableLoad, CommitTime: uint64(time.Now().UnixNano()), Schema: "reference_schema", Table: "reference_tables"},
		}))
		require.NoError(t, sinker.Push(reference.Table()))
		require.NoError(t, sinker.Push([]abstract.ChangeItem{
			{Kind: abstract.DoneTableLoad, CommitTime: uint64(time.Now().UnixNano()), Schema: "reference_schema", Table: "reference_tables"},
		}))
		source := &yt.YtSource{
			Cluster:          os.Getenv("YT_PROXY"),
			Proxy:            os.Getenv("YT_PROXY"),
			Paths:            []string{Destination.Path},
			YtToken:          "",
			RowIdxColumnName: "row_idx",
		}
		source.WithDefaults()
		reference.Canon(t, source)
	})
}
