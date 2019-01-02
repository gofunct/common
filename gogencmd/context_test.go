package gogencmd_test

import (
	"github.com/gofunct/common/files"
	"testing"

	"github.com/spf13/afero"

	"github.com/gofunct/gogen/pkg/cli"
	"github.com/gofunct/gogen/pkg/gogencmd"
)

func TestCtx(t *testing.T) {
	root := cli.RootDir{files.Path("/go/src/awesomeapp")}
	cwd := root.Join("api").String()

	orDie := func(t *testing.T, err error) {
		t.Helper()
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}
	}

	fs := afero.NewMemMapFs()
	orDie(t, fs.MkdirAll(cwd, 0755))
	orDie(t, afero.WriteFile(fs, root.Join("gogen.toml").String(), []byte(`
package = "awesomeapp"

[gogen]
server_dir = "./app/server"

[protoc]
protos_dir = "./api/protos"
out_dir = "./api"
import_dirs = [
  "./api/protos",
  "./vendor/github.com/grpc-ecosystem/grpc-gateway",
  "./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis",
]

  [[protoc.plugins]]
  name = "go"
  args = { plugins = "grpc", paths = "source_relative" }

  [[protoc.plugins]]
  name = "grpc-gateway"
  args = { logtostderr = true, paths = "source_relative" }

  [[protoc.plugins]]
  name = "swagger"
  args = { logtostderr = true }
`), 0644))

	ctx := &gogencmd.Ctx{FS: fs, RootDir: cli.RootDir{files.Path(cwd)}}

	err := ctx.Init()

	if err != nil {
		t.Errorf("Init() returned %v", err)
	}

	if got, want := ctx.RootDir, root; got != want {
		t.Errorf("RootDir is %q, want %q", got, want)
	}

	if got, want := ctx.IsInsideApp(), true; got != want {
		t.Errorf("IsInsideApp() returned %t, want %t", got, want)
	}

	if got, want := ctx.Config.Package, "awesomeapp"; got != want {
		t.Errorf("Config.Package is %q, want %q", got, want)
	}

	if got, want := ctx.ProtocConfig.ProtosDir, "./api/protos"; got != want {
		t.Errorf("ProtocConfig.ProtosDir is %q, want %q", got, want)
	}

	if got, want := len(ctx.ProtocConfig.Plugins), 3; got != want {
		t.Errorf("ProtocConfig has %d plugins, want %d", got, want)
	}
}

func TestCtx_outsideApp(t *testing.T) {
	root := cli.RootDir{files.Path("/go/src/awesomeapp")}
	cwd := root.Join("api").String()

	orDie := func(t *testing.T, err error) {
		t.Helper()
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}
	}

	fs := afero.NewMemMapFs()
	orDie(t, fs.MkdirAll(cwd, 0755))

	ctx := &gogencmd.Ctx{FS: fs, RootDir: cli.RootDir{files.Path(cwd)}}

	err := ctx.Init()

	if err != nil {
		t.Errorf("Init() returned %v", err)
	}

	if got, want := ctx.RootDir.String(), cwd; got != want {
		t.Errorf("RootDir is %q, want %q", got, want)
	}

	if got, want := ctx.IsInsideApp(), false; got != want {
		t.Errorf("IsInsideApp() returned %t, want %t", got, want)
	}

	if got, want := ctx.Config.Package, ""; got != want {
		t.Errorf("Config.Package is %q, want %q", got, want)
	}

	if got, want := ctx.ProtocConfig.ProtosDir, ""; got != want {
		t.Errorf("ProtocConfig.ProtosDir is %q, want %q", got, want)
	}

	if got, want := len(ctx.ProtocConfig.Plugins), 0; got != want {
		t.Errorf("ProtocConfig has %d plugins, want %d", got, want)
	}
}
