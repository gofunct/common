package exec

import (
	"fmt"
	"github.com/google/wire"
	"github.com/spf13/cobra"
	"strings"
)

// New returns a new Interface which will os/exec to run commands.
func NewCommander(s *Scripter) *cobra.Command {
	c := &cobra.Command{
		Use:         s.Name,
		Short:       s.Usage,
		Annotations: envtoAnnotations(s.OS.Env),
	}
	if s.RunFunc != nil {
		c.RunE = func(cmd *cobra.Command, args []string) error {
			return s.Execute(args...)
		}
	}
	c.SetOutput(s.IO.OutW)
	return c
}

var DefaultSet = wire.NewSet(
	NewCommander,
)

func envtoAnnotations(e []string) map[string]string {
	var a = make(map[string]string)

	for _, str := range e {

		sli := strings.Split(str, "=")
		a[sli[0]] = sli[1]
	}
	for k, v := range a {
		fmt.Println(k, v)
	}
	return a
}
