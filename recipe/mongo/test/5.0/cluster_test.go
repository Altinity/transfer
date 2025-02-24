package example

import (
	"testing"

	"github.com/altinity/transfer/recipe/mongo/pkg/util"
)

func TestSample(t *testing.T) {
	util.TestMongoShardedClusterRecipe(t)
}
