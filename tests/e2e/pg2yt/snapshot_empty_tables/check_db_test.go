package snapshot

import (
	"os"
	"testing"

	"github.com/altinity/transfer/pkg/abstract"
	"github.com/altinity/transfer/pkg/abstract/model"
	"github.com/altinity/transfer/pkg/providers/postgres"
	yt_provider "github.com/altinity/transfer/pkg/providers/yt"
	"github.com/altinity/transfer/tests/helpers"
	yt_helpers "github.com/altinity/transfer/tests/helpers/yt"
	"github.com/stretchr/testify/require"
)

var (
	Source = postgres.PgSource{
		ClusterID: os.Getenv("PG_CLUSTER_ID"),
		Hosts:     []string{"localhost"},
		User:      os.Getenv("PG_LOCAL_USER"),
		Password:  model.SecretString(os.Getenv("PG_LOCAL_PASSWORD")),
		Database:  os.Getenv("PG_LOCAL_DATABASE"),
		Port:      helpers.GetIntFromEnv("PG_LOCAL_PORT"),
		DBTables:  []string{"public.__test_empty"},
	}
	Target = yt_helpers.RecipeYtTarget("//home/cdc/test/pg2yt_e2e").(*yt_provider.YtDestinationWrapper)
)

func init() {
	_ = os.Setenv("YC", "1") // to not go to vanga
	Source.WithDefaults()
}

func TestMain(m *testing.M) {
	yt_provider.InitExe()
	os.Exit(m.Run())
}

func TestGroup(t *testing.T) {
	targetPort, err := helpers.GetPortFromStr(Target.Cluster())
	require.NoError(t, err)
	defer func() {
		require.NoError(t, helpers.CheckConnections(
			helpers.LabeledPort{Label: "PG source", Port: Source.Port},
			helpers.LabeledPort{Label: "YT target", Port: targetPort},
		))
	}()

	t.Run("Group after port check", func(t *testing.T) {
		t.Run("Snapshot", Snapshot)
	})
}

func Snapshot(t *testing.T) {
	Source.PreSteps.Constraint = true
	Target.Model.UseStaticTableOnSnapshot = true
	Target.Model.CreateEmptyTables = true
	transfer := helpers.MakeTransfer(helpers.TransferID, &Source, Target, abstract.TransferTypeSnapshotOnly)

	_ = helpers.Activate(t, transfer)

	require.NoError(t, helpers.CompareStorages(t, Source, Target.LegacyModel(), helpers.NewCompareStorageParams()))
}
