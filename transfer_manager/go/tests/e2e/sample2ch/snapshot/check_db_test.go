package replication

import (
	"testing"

	"github.com/doublecloud/transfer/transfer_manager/go/pkg/abstract"
	server "github.com/doublecloud/transfer/transfer_manager/go/pkg/abstract/model"
	chrecipe "github.com/doublecloud/transfer/transfer_manager/go/pkg/providers/clickhouse/recipe"
	"github.com/doublecloud/transfer/transfer_manager/go/pkg/providers/sample"
	"github.com/doublecloud/transfer/transfer_manager/go/tests/helpers"
	"github.com/stretchr/testify/require"
)

const expectedNumberOfRows = 100

var (
	schemaName   = "mtmobproxy"
	TransferType = abstract.TransferTypeSnapshotOnly
	Source       = *sample.RecipeSource()
	Target       = *chrecipe.MustTarget(chrecipe.WithInitFile("dump/dst.sql"), chrecipe.WithDatabase(schemaName), chrecipe.WithPrefix("DB0_"))
)

func TestSnapshot(t *testing.T) {
	defer func() {
		require.NoError(t, helpers.CheckConnections(
			helpers.LabeledPort{Label: "CH target", Port: Target.NativePort},
		))
	}()
	Target.WithDefaults()
	Target.Cleanup = server.DisabledCleanup

	Source.WithDefaults()
	helpers.InitSrcDst(helpers.TransferID, &Source, &Target, TransferType)

	transfer := helpers.MakeTransfer(helpers.TransferID, &Source, &Target, TransferType)

	helpers.Activate(t, transfer)

	helpers.CheckRowsCount(t, &Target, schemaName, "iot", expectedNumberOfRows)
}
