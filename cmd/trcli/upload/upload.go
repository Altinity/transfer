package upload

import (
	"context"

	"github.com/altinity/transfer/cmd/trcli/config"
	"github.com/altinity/transfer/library/go/core/metrics"
	"github.com/altinity/transfer/library/go/core/xerrors"
	"github.com/altinity/transfer/pkg/abstract"
	"github.com/altinity/transfer/pkg/abstract/coordinator"
	"github.com/altinity/transfer/pkg/abstract/model"
	"github.com/altinity/transfer/pkg/worker/tasks"
	"github.com/spf13/cobra"
)

func UploadCommand(cp *coordinator.Coordinator, rt abstract.Runtime, registry metrics.Registry) *cobra.Command {
	var transferParams string
	var uploadParams string
	var metricsPrefix string

	uploadCommand := &cobra.Command{
		Use:     "upload",
		Short:   "Upload tables",
		Example: "./trcli upload --transfer ./transfer.yaml --tables tables.yaml",
		RunE:    upload(cp, rt, &transferParams, &uploadParams, registry, metricsPrefix),
	}
	uploadCommand.Flags().StringVar(&transferParams, "transfer", "./transfer.yaml", "path to yaml file with transfer configuration")
	uploadCommand.Flags().StringVar(&uploadParams, "tables", "./tables.yaml", "path to yaml file with uploadable table params")
	uploadCommand.Flags().StringVar(&metricsPrefix, "metrics-prefix", "", "Optional prefix por Prometheus metrics")

	return uploadCommand
}

func upload(cp *coordinator.Coordinator, rt abstract.Runtime, transferYaml, uploadTablesYaml *string, registry metrics.Registry, metricsPrefix string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		transfer, err := config.TransferFromYaml(transferYaml)
		if err != nil {
			return xerrors.Errorf("unable to load transfer: %w", err)
		}

		transfer.Runtime = rt

		tables, err := config.TablesFromYaml(uploadTablesYaml)
		if err != nil {
			return xerrors.Errorf("unable to load tables: %w", err)
		}

		if metricsPrefix != "" {
			registry = registry.WithPrefix(metricsPrefix)
		}

		return RunUpload(*cp, transfer, tables, registry)
	}
}

func RunUpload(cp coordinator.Coordinator, transfer *model.Transfer, tables *config.UploadTables, registry metrics.Registry) error {
	return tasks.Upload(
		context.Background(),
		cp,
		*transfer,
		nil,
		tasks.UploadSpec{Tables: tables.Tables},
		registry.WithTags(map[string]string{
			"resource_id": transfer.ID,
			"name":        transfer.TransferName,
		}),
	)
}
