package executor_test

import (
	"github.com/gofunct/common/executor"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestOption_WithPATH(t *testing.T) {
	cases := []struct {
		test string
		env  []string
		in   string
		out  []string
	}{
		{
			test: "set path",
			env:  []string{"FOO=1", "BAR=baz"},
			in:   "/home/go/src/awsomeapp/bin",
			out:  []string{"FOO=1", "BAR=baz", "PATH=/home/go/src/awsomeapp/bin"},
		},
		{
			test: "append path",
			env:  []string{"FOO=1", "PATH=/home/go/bin", "BAR=baz"},
			in:   "/home/go/src/awsomeapp/bin",
			out:  []string{"FOO=1", "PATH=/home/go/src/awsomeapp/bin" + string(filepath.ListSeparator) + "/home/go/bin", "BAR=baz"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.test, func(t *testing.T) {
			cmd := &executor.Command{Env: tc.env}
			executor.WithPATH(tc.in)(cmd)

			if diff := cmp.Diff(cmd.Env, tc.out); diff != "" {
				t.Errorf("after WithPath: (-want +got): %s", diff)
			}
		})
	}
}
