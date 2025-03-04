package tests

import (
	_ "embed"
	"testing"
	"time"

	"github.com/altinity/transfer/cmd/trcli/activate"
	"github.com/altinity/transfer/cmd/trcli/config"
	"github.com/altinity/transfer/library/go/core/metrics/solomon"
	"github.com/altinity/transfer/pkg/abstract/coordinator"
	chrecipe "github.com/altinity/transfer/pkg/providers/clickhouse/recipe"
	"github.com/altinity/transfer/pkg/providers/postgres/pgrecipe"
	"github.com/altinity/transfer/tests/helpers"
	"github.com/stretchr/testify/require"
)

//go:embed transfer.yaml
var transferYaml []byte

func TestActivate(t *testing.T) {
	src := pgrecipe.RecipeSource(
		pgrecipe.WithPrefix(""),
		pgrecipe.WithFiles("dump/pg_init.sql"),
	)

	dst, err := chrecipe.Target(
		chrecipe.WithInitFile("ch_init.sql"),
		chrecipe.WithDatabase("trcli_activate_test_ch"),
	)
	require.NoError(t, err)

	transfer, err := config.ParseTransfer(transferYaml)
	require.NoError(t, err)

	transfer.Src = src
	transfer.Dst = dst

	require.NoError(t, activate.RunActivate(coordinator.NewStatefulFakeClient(), transfer, solomon.NewRegistry(solomon.NewRegistryOpts()), 0))

	require.NoError(t, helpers.WaitDestinationEqualRowsCount(dst.Database, "t2", helpers.GetSampleableStorageByModel(t, dst), 60*time.Second, 2))
	require.NoError(t, helpers.WaitDestinationEqualRowsCount(dst.Database, "t3", helpers.GetSampleableStorageByModel(t, dst), 60*time.Second, 5))
}

func TestActivateWithDelay(t *testing.T) {
	src := pgrecipe.RecipeSource(
		pgrecipe.WithPrefix(""),
		pgrecipe.WithFiles("dump/pg_init.sql"),
	)

	dst, err := chrecipe.Target(
		chrecipe.WithInitFile("ch_init.sql"),
		chrecipe.WithDatabase("trcli_activate_test_ch"),
	)
	require.NoError(t, err)

	transfer, err := config.ParseTransfer(transferYaml)
	require.NoError(t, err)

	transfer.Src = src
	transfer.Dst = dst

	st := time.Now()
	require.NoError(t, activate.RunActivate(coordinator.NewStatefulFakeClient(), transfer, solomon.NewRegistry(solomon.NewRegistryOpts()), 10*time.Second))
	require.Less(t, 10*time.Second, time.Since(st))
}
