package transformer_test

import (
	"testing"

	"github.com/altinity/transfer/internal/logger"
	"github.com/altinity/transfer/library/go/core/metrics/solomon"
	"github.com/altinity/transfer/pkg/abstract"
	"github.com/altinity/transfer/pkg/abstract/coordinator"
	"github.com/altinity/transfer/pkg/abstract/model"
	"github.com/altinity/transfer/pkg/middlewares"
	"github.com/altinity/transfer/pkg/sink"
	transformers_registry "github.com/altinity/transfer/pkg/transformer"
	"github.com/altinity/transfer/pkg/transformer/registry/filter"
	replaceprimarykey "github.com/altinity/transfer/pkg/transformer/registry/replace_primary_key"
	"github.com/stretchr/testify/require"
)

type mockSinker struct {
	gotItems []abstract.ChangeItem
}

func (m *mockSinker) Close() error { return nil }
func (m *mockSinker) Push(items []abstract.ChangeItem) error {
	m.gotItems = append(m.gotItems, items...)
	return nil
}

func TestMultipleTransformers(t *testing.T) {
	tableName := "test_table"
	trans := &model.Transformation{
		Transformers: &transformers_registry.Transformers{
			DebugMode: true,
			Transformers: []transformers_registry.Transformer{
				{
					replaceprimarykey.Type: replaceprimarykey.Config{
						Keys: []string{
							"field2",
							"field1",
						},
						Tables: filter.Tables{
							IncludeTables: []string{tableName},
						},
					},
				},
				{
					filter.FilterColumnsTransformerType: filter.FilterColumnsConfig{
						Tables: filter.Tables{
							IncludeTables: []string{tableName},
						},
						Columns: filter.Columns{
							IncludeColumns: []string{
								"field2",
								"field1",
								"field4",
							},
						},
					},
				},
			},
			ErrorsOutput: nil,
		},
		ExtraTransformers: nil,
		Executor:          nil,
	}

	mockSinker := new(mockSinker)
	transfer := &model.Transfer{
		Src: &model.MockSource{},
		Dst: &model.MockDestination{
			SinkerFactory: func() abstract.Sinker { return mockSinker },
			Cleanup:       model.Drop,
		},
		Transformation: trans,
	}
	asink, err := sink.MakeAsyncSink(
		transfer,
		logger.Log,
		solomon.NewRegistry(solomon.NewRegistryOpts()),
		coordinator.NewFakeClient(),
		middlewares.MakeConfig(middlewares.WithNoData),
	)
	require.NoError(t, err)
	var data = []abstract.ChangeItem{
		abstract.ChangeItemFromMap(nil, nil, tableName, string(abstract.InitTableLoad)),
		abstract.ChangeItemFromMap(map[string]interface{}{
			"field1": "test",
			"field2": 2,
			"field3": 1.23,
			"field4": "{}",
		}, abstract.NewTableSchema([]abstract.ColSchema{
			{ColumnName: "field1", PrimaryKey: true},
			{ColumnName: "field2", PrimaryKey: true},
			{ColumnName: "field3", PrimaryKey: true},
			{ColumnName: "field4", PrimaryKey: true},
		}), tableName, string(abstract.InsertKind)),
	}

	require.NoError(t, <-asink.AsyncPush(data))
	require.NoError(t, asink.Close())
	require.Equal(t, len(mockSinker.gotItems), 2)
	require.Equal(t, mockSinker.gotItems[1].Kind, abstract.InsertKind)
	require.Equal(t, mockSinker.gotItems[1].TableSchema,
		abstract.NewTableSchema([]abstract.ColSchema{
			{ColumnName: "field2", PrimaryKey: true},
			{ColumnName: "field1", PrimaryKey: true},
			{ColumnName: "field4", PrimaryKey: false},
		}), tableName, string(abstract.InsertKind))
	require.Equal(t, mockSinker.gotItems[1].ColumnValues, []interface{}{"test", 2, "{}"})

}
