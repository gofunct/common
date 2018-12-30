package cmd

import (
	"github.com/gofunct/common/build"
	"github.com/gofunct/common/cli"
	"github.com/gofunct/common/files"
	"github.com/gofunct/common/io"
	"github.com/gofunct/common/logging"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"k8s.io/utils/exec"
)

func NewDefaultCligCommand(wd files.Path, build build.Build) *cobra.Command {
	return NewCligCommand(&cli.Ctx{
		WorkingDir: wd,
		IO:         io.Stdio(),
		FS:         afero.NewOsFs(),
		Exec:       exec.New(),
		Build:      build,
	})
}

func NewCligCommand(ctx *cli.Ctx) *cobra.Command {
	cmd := &cobra.Command{
		Use: ctx.Build.AppName,
	}

	logging.AddLoggingFlags(cmd)

	cmd.AddCommand(
		newInitCommand(ctx),
		build.NewVersionCommand(ctx.IO, ctx.Build),
	)

	return cmd
}
