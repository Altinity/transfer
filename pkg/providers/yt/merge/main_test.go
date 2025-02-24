package merge

import (
	"os"
	"testing"

	ytcommon "github.com/altinity/transfer/pkg/providers/yt"
)

func TestMain(m *testing.M) {
	ytcommon.InitExe()
	os.Exit(m.Run())
}
