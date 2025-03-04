package tests

import (
	_ "embed"
	"testing"
	"time"

	"github.com/altinity/transfer/cmd/trcli/config"
	"github.com/altinity/transfer/cmd/trcli/upload"
	"github.com/altinity/transfer/library/go/core/metrics/solomon"
	"github.com/altinity/transfer/pkg/abstract/coordinator"
	chrecipe "github.com/altinity/transfer/pkg/providers/clickhouse/recipe"
	"github.com/altinity/transfer/pkg/providers/postgres/pgrecipe"
	"github.com/altinity/transfer/tests/helpers"
	"github.com/stretchr/testify/require"
)

//go:embed transfer.yaml
var transferYaml []byte

//go:embed tables.yaml
var tablesYaml []byte

func TestUpload(t *testing.T) {
	src := pgrecipe.RecipeSource(
		pgrecipe.WithPrefix(""),
		pgrecipe.WithFiles("dump/pg_init.sql"),
	)

	dst, err := chrecipe.Target(
		chrecipe.WithInitFile("ch_init.sql"),
		chrecipe.WithDatabase("trcli_upload_test_ch"),
	)
	require.NoError(t, err)

	transfer, err := config.ParseTransfer(transferYaml)
	require.NoError(t, err)

	transfer.Src = src
	transfer.Dst = dst

	tables, err := config.ParseTablesYaml(tablesYaml)
	require.NoError(t, err)

	require.NoError(t, upload.RunUpload(coordinator.NewFakeClient(), transfer, tables, solomon.NewRegistry(solomon.NewRegistryOpts())))
	require.NoError(t, helpers.WaitDestinationEqualRowsCount(dst.Database, "t2", helpers.GetSampleableStorageByModel(t, dst), 60*time.Second, 2))
	require.NoError(t, helpers.WaitDestinationEqualRowsCount(dst.Database, "t3", helpers.GetSampleableStorageByModel(t, dst), 60*time.Second, 2))
}
