package cascadedeletespertbl

import (
	"testing"

	test "github.com/altinity/transfer/tests/e2e/mysql2mysql/cascade_deletes/common"
	"github.com/altinity/transfer/tests/helpers"
	"github.com/stretchr/testify/require"
)

func TestGroup(t *testing.T) {
	defer func() {
		require.NoError(t, helpers.CheckConnections(
			helpers.LabeledPort{Label: "Mysql source", Port: test.Source.Port},
			helpers.LabeledPort{Label: "Mysql target", Port: test.Target.Port},
		))
	}()

	t.Run("Group after port check", func(t *testing.T) {
		t.Run("Existence", test.Existence)
		t.Run("Snapshot", test.Snapshot)
		t.Run("Replication", test.Load)
	})
}
