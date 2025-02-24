package dataplane

import (
	_ "github.com/altinity/transfer/pkg/providers/airbyte"
	_ "github.com/altinity/transfer/pkg/providers/clickhouse"
	_ "github.com/altinity/transfer/pkg/providers/coralogix"
	_ "github.com/altinity/transfer/pkg/providers/datadog"
	_ "github.com/altinity/transfer/pkg/providers/delta"
	_ "github.com/altinity/transfer/pkg/providers/elastic"
	_ "github.com/altinity/transfer/pkg/providers/eventhub"
	_ "github.com/altinity/transfer/pkg/providers/greenplum"
	_ "github.com/altinity/transfer/pkg/providers/kafka"
	_ "github.com/altinity/transfer/pkg/providers/mongo"
	_ "github.com/altinity/transfer/pkg/providers/mysql"
	_ "github.com/altinity/transfer/pkg/providers/opensearch"
	_ "github.com/altinity/transfer/pkg/providers/postgres"
	_ "github.com/altinity/transfer/pkg/providers/s3/provider"
	_ "github.com/altinity/transfer/pkg/providers/stdout"
	_ "github.com/altinity/transfer/pkg/providers/ydb"
	_ "github.com/altinity/transfer/pkg/providers/yt/init"
)
