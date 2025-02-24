package prodstatus

import (
	"github.com/altinity/transfer/pkg/abstract"
	"github.com/altinity/transfer/pkg/providers/mysql"
	"github.com/altinity/transfer/pkg/providers/postgres"
	"github.com/altinity/transfer/pkg/providers/ydb"
)

var supportedSources = map[string]bool{
	postgres.ProviderType.Name(): true,
	mysql.ProviderType.Name():    true,
	ydb.ProviderType.Name():      true,
}

func IsSupportedSource(src string, _ abstract.TransferType) bool {
	return supportedSources[src]
}
