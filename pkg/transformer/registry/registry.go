package registry

import (
	_ "github.com/altinity/transfer/pkg/transformer/registry/batch_splitter"
	_ "github.com/altinity/transfer/pkg/transformer/registry/clickhouse"
	_ "github.com/altinity/transfer/pkg/transformer/registry/custom"
	_ "github.com/altinity/transfer/pkg/transformer/registry/filter"
	_ "github.com/altinity/transfer/pkg/transformer/registry/filter_rows"
	_ "github.com/altinity/transfer/pkg/transformer/registry/logger"
	_ "github.com/altinity/transfer/pkg/transformer/registry/mask"
	_ "github.com/altinity/transfer/pkg/transformer/registry/number_to_float"
	_ "github.com/altinity/transfer/pkg/transformer/registry/problem_item_detector"
	_ "github.com/altinity/transfer/pkg/transformer/registry/raw_doc_grouper"
	_ "github.com/altinity/transfer/pkg/transformer/registry/rename"
	_ "github.com/altinity/transfer/pkg/transformer/registry/replace_primary_key"
	_ "github.com/altinity/transfer/pkg/transformer/registry/sharder"
	_ "github.com/altinity/transfer/pkg/transformer/registry/table_splitter"
	_ "github.com/altinity/transfer/pkg/transformer/registry/to_string"
	_ "github.com/altinity/transfer/pkg/transformer/registry/yt_dict"
)
