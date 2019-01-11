package exec

import (
	"context"
	"github.com/spf13/cobra"
	"os/exec"
)

type Interface interface {
	Execute() error
}

type ScriptFunc func(exec.Cmd) Script

type CobraFunc func(*cobra.Command) Script

type Script func(ctx context.Context, name string, args ...string) error
