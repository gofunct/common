package logging_test

import (
	"github.com/gofunct/common/io"
	"github.com/spf13/cobra"
	"strings"
	"testing"
	"github.com/gofunct/common/logging"
)

func TestEmojies(t *testing.T) {
	logger := logging.NewLogger()

	logger.UI.Output("output")
	logger.UI.Running("running")
	logger.UI.Success("success")
	logger.UI.Log("log")
	logger.UI.Error("error")
	logger.UI.Info("info")
	logger.UI.Ask("?", "ask")
	logger.UI.Warn("warn")
}


func TestLogging(t *testing.T) {
	cases := []struct {
		args      []string
		mode      logging.LoggingMode
		isDebug   bool
		isVerbose bool
	}{
		{
			mode: logging.LoggingNop,
		},
		{
			args:      []string{"-v"},
			mode:      logging.LoggingVerbose,
			isVerbose: true,
		},
		{
			args:      []string{"--verbose"},
			mode:      logging.LoggingVerbose,
			isVerbose: true,
		},
		{
			args:    []string{"--debug"},
			mode:    logging.LoggingDebug,
			isDebug: true,
		},
	}

	for _, tc := range cases {
		t.Run(strings.Join(tc.args, " "), func(t *testing.T) {
			defer io.Close()

			var (
				mode               logging.LoggingMode
				isDebug, isVerbose bool
			)

			cmd := &cobra.Command{
				Run: func(*cobra.Command, []string) {
					mode = logging.Logging()
					isDebug = logging.IsDebug()
					isVerbose = logging.IsVerbose()
				},
			}

			logging.AddLoggingFlags(cmd)
			cmd.SetArgs(tc.args)
			err := cmd.Execute()

			if err != nil {
				t.Errorf("Execute() returned an error: %v", err)
			}

			if got, want := mode, tc.mode; got != want {
				t.Errorf("LoggingMode() returned %v, want %v", got, want)
			}

			if got, want := isVerbose, tc.isVerbose; got != want {
				t.Errorf("IsVerbose() returned %t, want %t", got, want)
			}

			if got, want := isDebug, tc.isDebug; got != want {
				t.Errorf("IsDebug() returned %t, want %t", got, want)
			}
		})
	}
}
