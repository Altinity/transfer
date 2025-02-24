package tests

import (
	"context"
	_ "embed"
	"testing"
	"time"

	"github.com/altinity/transfer/cmd/trcli/config"
	"github.com/altinity/transfer/cmd/trcli/replicate"
	"github.com/altinity/transfer/internal/logger"
	"github.com/altinity/transfer/library/go/core/metrics/solomon"
	"github.com/altinity/transfer/pkg/abstract/coordinator"
	chrecipe "github.com/altinity/transfer/pkg/providers/clickhouse/recipe"
	pgcommon "github.com/altinity/transfer/pkg/providers/postgres"
	"github.com/altinity/transfer/pkg/providers/postgres/pgrecipe"
	"github.com/altinity/transfer/tests/helpers"
	"github.com/stretchr/testify/require"
)

//go:embed transfer.yaml
var transferYaml []byte

func TestReplicate(t *testing.T) {
	src := pgrecipe.RecipeSource(
		pgrecipe.WithPrefix(""),
		pgrecipe.WithFiles("dump/pg_init.sql"),
	)

	dst, err := chrecipe.Target(
		chrecipe.WithInitFile("ch_init.sql"),
		chrecipe.WithDatabase("trcli_replicate_test_ch"),
	)
	require.NoError(t, err)
	transfer, err := config.ParseTransfer(transferYaml)
	require.NoError(t, err)

	src.SlotID = transfer.ID
	transfer.Src = src
	transfer.Dst = dst

	go func() {
		require.NoError(t, replicate.RunReplication(coordinator.NewStatefulFakeClient(), transfer, solomon.NewRegistry(solomon.NewRegistryOpts())))
	}()

	time.Sleep(5 * time.Second)

	connConfig, err := pgcommon.MakeConnConfigFromSrc(logger.Log, src)
	require.NoError(t, err)

	conn, err := pgcommon.NewPgConnPool(connConfig, logger.Log)
	require.NoError(t, err)

	rows, err := conn.Query(context.Background(), "INSERT INTO public.t2(i, f) VALUES (3, 1.0), (4, 4.0)")
	require.NoError(t, err)
	rows.Close()

	rows, err = conn.Query(context.Background(), "INSERT INTO public.t3(i, f) VALUES (1, 2.0), (2, 3.0)")
	require.NoError(t, err)
	rows.Close()

	require.NoError(t, helpers.WaitDestinationEqualRowsCount(dst.Database, "t2", helpers.GetSampleableStorageByModel(t, dst), 60*time.Second, 2))
	require.NoError(t, helpers.WaitDestinationEqualRowsCount(dst.Database, "t3", helpers.GetSampleableStorageByModel(t, dst), 60*time.Second, 2))
}
