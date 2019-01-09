package exec

import (
	"context"
	"github.com/gofunct/iio"
	"github.com/spf13/cobra"
	"os/exec"
)

//TODO:Refine Scripter
type StdScripter struct {
	Name         string
	Initializers []func()
	c            *cobra.Command
	e            *exec.Cmd
}

func (s *StdScripter) OnInitialize(f ...func()) {
	cobra.OnInitialize(f...)
}

func (s *StdScripter) Execute(ctx context.Context, i iio.IO) ([]byte, error) {
	var out []byte

	cob := &cobra.Command{
		Use:   s.Name,
		Short: s.e.Path,
		PreRun: func(cmd *cobra.Command, args []string) {
			s.OnInitialize(s.Initializers...)
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			cmd.SetOutput(i.Out)
			var err error
			s.e = exec.CommandContext(ctx, s.Name, args...)
			s.e.Stderr = i.Err
			s.e.Stdout = i.Out
			s.e.Stdin = i.In
			out, err = s.e.Output()
			return err
		},
	}

	return out, cob.Execute()
}
