package case1

import (
	"os"
	"testing"
	"time"

	"github.com/altinity/transfer/pkg/abstract"
	"github.com/altinity/transfer/pkg/abstract/model"
	ytcommon "github.com/altinity/transfer/pkg/providers/yt"
	"github.com/altinity/transfer/tests/e2e/mongo2yt/rotator"
	"go.ytsaurus.tech/yt/go/ypath"
)

func TestMain(m *testing.M) {
	ytcommon.InitExe()
	os.Exit(m.Run())
}

func TestCases(t *testing.T) {

	// fix time with modern but certain point
	// Note that rotator may delete tables if date is too far away, so 'now' value is strongly recommended
	ts := time.Now()

	table := abstract.TableID{Namespace: "db", Name: "test"}
	dayRotationExpectedTable := rotator.DayRotation.AnnotateWithTime("db_test", ts)

	t.Run("cleanup=drop;rotation=day;use_static_table=true;table_type=static", func(t *testing.T) {
		source, target := rotator.PrefilledSourceAndTarget()
		target.Cleanup = model.Drop
		target.Rotation = rotator.DayRotation
		target.UseStaticTableOnSnapshot = true
		target.Static = true
		expectedPath := ypath.Path(target.Path).Child(dayRotationExpectedTable)
		rotator.ScenarioCheckActivation(t, source, target, table, ts, expectedPath)
	})
}
