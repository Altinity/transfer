package events

import (
	"github.com/altinity/transfer/library/go/core/xerrors"
	"github.com/altinity/transfer/pkg/base"
)

func validateValue(value base.Value) error {
	if err := value.Column().Type().Validate(value); err != nil {
		return xerrors.Errorf("Column '%v', value validation error: %w", value.Column().FullName(), err)
	}
	return nil
}
