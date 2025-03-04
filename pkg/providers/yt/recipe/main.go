package recipe

import (
	"context"
	"os"
	"testing"

	ytcommon "github.com/altinity/transfer/pkg/providers/yt"
	"github.com/testcontainers/testcontainers-go"
)

func Main(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	container, _ := RunContainer(ctx, testcontainers.WithImage("ytsaurus/local:stable"))
	proxy, _ := container.ConnectionHost(ctx)
	_ = os.Setenv("YT_PROXY", proxy)
	ytcommon.InitExe()
	res := m.Run()
	_ = container.Terminate(ctx)
	cancel()
	os.Exit(res)
}
