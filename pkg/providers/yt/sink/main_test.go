package sink

import (
	"os"
	"testing"

	"github.com/altinity/transfer/pkg/config/env"
	ytcommon "github.com/altinity/transfer/pkg/providers/yt"
	"github.com/altinity/transfer/pkg/providers/yt/recipe"
)

func TestMain(m *testing.M) {
	if recipe.TestContainerEnabled() {
		recipe.Main(m)
		return
	}
	if env.IsTest() && !recipe.TestContainerEnabled() {
		ytcommon.InitExe()
	}
	os.Exit(m.Run())
}
