package replicationview

import (
	"os"
	"testing"
	"time"

	"github.com/altinity/transfer/internal/logger"
	"github.com/altinity/transfer/pkg/abstract"
	"github.com/altinity/transfer/pkg/abstract/coordinator"
	"github.com/altinity/transfer/pkg/abstract/model"
	"github.com/altinity/transfer/pkg/providers/postgres/pgrecipe"
	"github.com/altinity/transfer/pkg/runtime/local"
	"github.com/altinity/transfer/tests/helpers"
	"github.com/stretchr/testify/require"
)

func init() {
	_ = os.Setenv("YC", "1") // to not go to vanga
}

func TestViewReplication(t *testing.T) {
	Source := *pgrecipe.RecipeSource(pgrecipe.WithInitDir("init_source"))
	Target := *pgrecipe.RecipeTarget(pgrecipe.WithInitDir("init_target"))
	Target.Cleanup = model.Truncate

	defer func() {
		require.NoError(t, helpers.CheckConnections(
			helpers.LabeledPort{Label: "PG source", Port: Source.Port},
			helpers.LabeledPort{Label: "PG target", Port: Target.Port},
		))
	}()

	helpers.InitSrcDst(helpers.TransferID, &Source, &Target, abstract.TransferTypeIncrementOnly)
	worker := local.NewLocalWorker(coordinator.NewFakeClient(), helpers.MakeTransfer(helpers.TransferID, &Source, &Target, abstract.TransferTypeIncrementOnly), helpers.EmptyRegistry(), logger.Log)
	worker.Start()

	require.NoError(t, helpers.WaitEqualRowsCount(t, "public", "tv_table", helpers.GetSampleableStorageByModel(t, Source), helpers.GetSampleableStorageByModel(t, Target), 20*time.Second))
	require.NoError(t, helpers.WaitEqualRowsCount(t, "public", "odd_channels", helpers.GetSampleableStorageByModel(t, Source), helpers.GetSampleableStorageByModel(t, Target), 20*time.Second))

	require.NoError(t, worker.Stop())
}
