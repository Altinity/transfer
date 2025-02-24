package registry

import (
	_ "github.com/altinity/transfer/pkg/parsers/registry/audittrailsv1"
	_ "github.com/altinity/transfer/pkg/parsers/registry/blank"
	_ "github.com/altinity/transfer/pkg/parsers/registry/cloudevents"
	_ "github.com/altinity/transfer/pkg/parsers/registry/cloudlogging"
	_ "github.com/altinity/transfer/pkg/parsers/registry/confluentschemaregistry"
	_ "github.com/altinity/transfer/pkg/parsers/registry/debezium"
	_ "github.com/altinity/transfer/pkg/parsers/registry/json"
	_ "github.com/altinity/transfer/pkg/parsers/registry/native"
	_ "github.com/altinity/transfer/pkg/parsers/registry/protobuf"
	_ "github.com/altinity/transfer/pkg/parsers/registry/raw2table"
	_ "github.com/altinity/transfer/pkg/parsers/registry/tskv"
)
