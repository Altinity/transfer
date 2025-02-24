package fallback

import (
	"github.com/altinity/transfer/pkg/abstract"
	"github.com/altinity/transfer/pkg/abstract/typesystem"
	"github.com/altinity/transfer/pkg/providers/s3"
)

func init() {
	typesystem.AddFallbackTargetFactory(func() typesystem.Fallback {
		return typesystem.Fallback{
			To:     8,
			Picker: typesystem.ProviderType(s3.ProviderType),
			Function: func(ci *abstract.ChangeItem) (*abstract.ChangeItem, error) {
				if ci.Schema == "" {
					ci.Table = "_" + ci.Table
				}
				return ci, nil
			},
		}
	})
}
