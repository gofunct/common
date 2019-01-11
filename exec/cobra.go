package exec

import (
	"context"
	"github.com/spf13/cobra"
)

type Commander struct {
	*cobra.Command
	Run CobraFunc
}

func test(c *cobra.Command) CobraFunc {

	return func(c *cobra.Command) Script {

		return func(ctx context.Context, name string, args ...string) error {

			return nil
		}
	}
}
