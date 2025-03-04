package stdout

import (
	"github.com/altinity/transfer/pkg/abstract"
	"github.com/altinity/transfer/pkg/abstract/model"
)

type EmptySource struct{}

var _ model.Source = (*EmptySource)(nil)

func (EmptySource) WithDefaults() {
}

func (s *EmptySource) Include(tID abstract.TableID) bool {
	return len(s.FulfilledIncludes(tID)) > 0
}

func (*EmptySource) FulfilledIncludes(tID abstract.TableID) []string {
	return []string{""}
}

func (*EmptySource) AllIncludes() []string {
	return nil
}

func (EmptySource) IsSource() {
}

func (s *EmptySource) GetProviderType() abstract.ProviderType {
	return ProviderType
}

func (s *EmptySource) Validate() error {
	return nil
}
