package mongo

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/altinity/transfer/pkg/abstract/typesystem"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed typesystem.md
	canonDoc string
)

func TestTypeSystem(t *testing.T) {
	rules := typesystem.RuleFor(ProviderType)
	require.NotNil(t, rules.Source)
	require.NotNil(t, rules.Target)
	doc := typesystem.Doc(ProviderType, "MongoDB")
	fmt.Print(doc)
	require.Equal(t, canonDoc, doc)
}
