package canon

import (
	"testing"

	"github.com/altinity/transfer/pkg/providers/clickhouse"
	"github.com/altinity/transfer/pkg/providers/mongo"
	"github.com/altinity/transfer/pkg/providers/mysql"
	"github.com/altinity/transfer/pkg/providers/postgres"
	"github.com/altinity/transfer/pkg/providers/ydb"
	"github.com/altinity/transfer/pkg/providers/yt"
	"github.com/altinity/transfer/tests/canon/validator"
	"github.com/stretchr/testify/require"
)

func TestAll(t *testing.T) {
	cases := All(
		ydb.ProviderType,
		yt.ProviderType,
		mongo.ProviderType,
		clickhouse.ProviderType,
		mysql.ProviderType,
		postgres.ProviderType,
	)
	for _, tc := range cases {
		t.Run(tc.String(), func(t *testing.T) {
			require.NotEmpty(t, tc.Data)
			snkr := validator.Referencer(t)()
			require.NoError(t, snkr.Push(tc.Data))
			require.NoError(t, snkr.Close())
		})
	}
}
