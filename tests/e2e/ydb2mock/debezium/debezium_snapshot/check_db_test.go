package main

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/altinity/transfer/internal/logger"
	"github.com/altinity/transfer/library/go/core/metrics/solomon"
	"github.com/altinity/transfer/library/go/test/yatest"
	"github.com/altinity/transfer/pkg/abstract"
	"github.com/altinity/transfer/pkg/abstract/model"
	debeziumcommon "github.com/altinity/transfer/pkg/debezium/common"
	"github.com/altinity/transfer/pkg/debezium/testutil"
	"github.com/altinity/transfer/pkg/providers/ydb"
	"github.com/altinity/transfer/tests/helpers"
	"github.com/stretchr/testify/require"
)

//---------------------------------------------------------------------------------------------------------------------

func TestGroup(t *testing.T) {
	src := &ydb.YdbSource{
		Token:              model.SecretString(os.Getenv("YDB_TOKEN")),
		Database:           helpers.GetEnvOfFail(t, "YDB_DATABASE"),
		Instance:           helpers.GetEnvOfFail(t, "YDB_ENDPOINT"),
		Tables:             nil,
		TableColumnsFilter: nil,
		SubNetworkID:       "",
		Underlay:           false,
		ServiceAccountID:   "",
	}

	//-----------------------------------------------------------------------------------------------------------------

	canonizedDebeziumKeyArr, err := os.ReadFile(yatest.SourcePath("transfer_manager/go/tests/e2e/ydb2mock/debezium/debezium_snapshot/testdata/change_item_key.txt"))
	require.NoError(t, err)
	canonizedDebeziumValArr, err := os.ReadFile(yatest.SourcePath("transfer_manager/go/tests/e2e/ydb2mock/debezium/debezium_snapshot/testdata/change_item_val.txt"))
	require.NoError(t, err)
	canonizedDebeziumVal := string(canonizedDebeziumValArr)

	//-----------------------------------------------------------------------------------------------------------------
	// init

	t.Run("init source database", func(t *testing.T) {
		Target := &ydb.YdbDestination{
			Database: src.Database,
			Token:    src.Token,
			Instance: src.Instance,
		}
		Target.WithDefaults()
		sinker, err := ydb.NewSinker(logger.Log, Target, solomon.NewRegistry(solomon.NewRegistryOpts()))
		require.NoError(t, err)

		currChangeItem := helpers.YDBInitChangeItem("dectest/timmyb32r-test")
		require.NoError(t, sinker.Push([]abstract.ChangeItem{*currChangeItem}))
	})

	//-----------------------------------------------------------------------------------------------------------------
	// activate

	sinker := &helpers.MockSink{}
	target := model.MockDestination{
		SinkerFactory: func() abstract.Sinker { return sinker },
		Cleanup:       model.DisabledCleanup,
	}
	transfer := helpers.MakeTransfer("fake", src, &target, abstract.TransferTypeSnapshotOnly)

	var changeItems []abstract.ChangeItem
	sinker.PushCallback = func(input []abstract.ChangeItem) {
		changeItems = append(changeItems, input...)
	}

	helpers.Activate(t, transfer)

	//-----------------------------------------------------------------------------------------------------------------
	// check

	require.Equal(t, 5, len(changeItems))
	require.Equal(t, changeItems[0].Kind, abstract.InitShardedTableLoad)
	require.Equal(t, changeItems[1].Kind, abstract.InitTableLoad)
	require.Equal(t, changeItems[2].Kind, abstract.InsertKind)
	require.Equal(t, changeItems[3].Kind, abstract.DoneTableLoad)
	require.Equal(t, changeItems[4].Kind, abstract.DoneShardedTableLoad)

	logger.Log.Infof("changeItem dump: %s\n", changeItems[2].ToJSONString())

	testutil.CheckCanonizedDebeziumEvent(t, &changeItems[2], "fullfillment", "pguser", "pg", true, []debeziumcommon.KeyValue{{DebeziumKey: string(canonizedDebeziumKeyArr), DebeziumVal: &canonizedDebeziumVal}})
	changeItemBuf, err := json.Marshal(changeItems[2])
	require.NoError(t, err)
	changeItemDeserialized := helpers.UnmarshalChangeItem(t, changeItemBuf)
	testutil.CheckCanonizedDebeziumEvent(t, changeItemDeserialized, "fullfillment", "pguser", "pg", true, []debeziumcommon.KeyValue{{DebeziumKey: string(canonizedDebeziumKeyArr), DebeziumVal: &canonizedDebeziumVal}})
}
