package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/altinity/transfer/pkg/abstract"
	"github.com/altinity/transfer/pkg/abstract/coordinator"
	"github.com/altinity/transfer/pkg/abstract/model"
	"github.com/altinity/transfer/pkg/providers/postgres/pgrecipe"
	"github.com/altinity/transfer/pkg/worker/tasks"
	"github.com/altinity/transfer/tests/helpers"
	"github.com/stretchr/testify/require"
)

var Source = *pgrecipe.RecipeSource(pgrecipe.WithInitDir("init_source"))

func init() {
	_ = os.Setenv("YC", "1") // to not go to vanga
	Source.WithDefaults()
}

//---------------------------------------------------------------------------------------------------------------------
// mockSinker

type mockSinker struct {
	pushCallback func([]abstract.ChangeItem)
}

func (s *mockSinker) Close() error {
	return nil
}

func (s *mockSinker) Push(input []abstract.ChangeItem) error {
	s.pushCallback(input)
	return nil
}

//---------------------------------------------------------------------------------------------------------------------

func TestGroup(t *testing.T) {
	defer func() {
		require.NoError(t, helpers.CheckConnections(
			helpers.LabeledPort{Label: "PG source", Port: Source.Port},
		))
	}()

	t.Run("Group after port check", func(t *testing.T) {
		t.Run("Snapshot", Snapshot)
	})
}

func Snapshot(t *testing.T) {
	sinker := &mockSinker{}
	transfer := helpers.MakeTransfer(
		helpers.TransferID,
		&Source,
		&model.MockDestination{SinkerFactory: func() abstract.Sinker {
			return sinker
		}},
		abstract.TransferTypeSnapshotOnly,
	)

	inputs := make(chan []abstract.ChangeItem, 100)
	sinker.pushCallback = func(input []abstract.ChangeItem) {
		time.Sleep(6 * time.Second)
		inputs <- input
	}

	tables, err := tasks.ObtainAllSrcTables(transfer, helpers.EmptyRegistry())
	require.NoError(t, err)
	snapshotLoader := tasks.NewSnapshotLoader(coordinator.NewFakeClient(), "test-operation", transfer, helpers.EmptyRegistry())
	err = snapshotLoader.UploadTables(context.Background(), tables.ConvertToTableDescriptions(), true)
	require.NoError(t, err)
}
