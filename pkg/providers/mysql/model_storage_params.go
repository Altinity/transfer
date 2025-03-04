package mysql

import "github.com/altinity/transfer/pkg/abstract"

type MysqlStorageParams struct {
	ClusterID   string
	Host        string
	Port        int
	User        string
	Password    string
	Database    string
	TLS         bool
	CertPEMFile string

	UseFakePrimaryKey   bool
	DegreeOfParallelism int
	Timezone            string

	TableFilter        abstract.Includeable
	PreSteps           *MysqlDumpSteps
	ConsistentSnapshot bool
	RootCAFiles        []string
	ConnectionID       string
}
