package exec

import (
	"github.com/gofunct/iio"
	"github.com/google/wire"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/exec"
)

// New returns a new Interface which will os/exec to run commands.
func New(name string, i *iio.Service, args ...string) Interface {
	return &Scripter{
		Cmd: exec.Cmd{
			Path:   name,
			Args:   args,
			Env:    viper.GetStringSlice("env"),
			Dir:    os.Getenv("PWD"),
			Stdin:  i.In(),
			Stdout: i.Out(),
			Stderr: i.Err(),
		},
	}
}

var DefaultSet = wire.NewSet(
	New,
)

func NewCobraCommand() *cobra.Command {
	return &cobra.Command{
		Use:                        "",
		Aliases:                    nil,
		SuggestFor:                 nil,
		Short:                      "",
		Long:                       "",
		Example:                    "",
		ValidArgs:                  nil,
		Args:                       nil,
		ArgAliases:                 nil,
		BashCompletionFunction:     "",
		Deprecated:                 "",
		Hidden:                     false,
		Annotations:                nil,
		Version:                    "",
		PersistentPreRun:           nil,
		PersistentPreRunE:          nil,
		PreRun:                     nil,
		PreRunE:                    nil,
		Run:                        nil,
		RunE:                       nil,
		PostRun:                    nil,
		PostRunE:                   nil,
		PersistentPostRun:          nil,
		PersistentPostRunE:         nil,
		SilenceErrors:              false,
		SilenceUsage:               false,
		DisableFlagParsing:         false,
		DisableAutoGenTag:          false,
		DisableFlagsInUseLine:      false,
		DisableSuggestions:         false,
		SuggestionsMinimumDistance: 0,
		TraverseChildren:           false,
		FParseErrWhitelist:         cobra.FParseErrWhitelist{},
	}
}
