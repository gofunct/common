package files_test

import (
	"github.com/gofunct/common/files"
	"testing"
	)

func TestPath_String(t *testing.T) {
	pathStr := "/go/src/awesomeapp"
	path := files.Path(pathStr)

	if got, want := path.String(), pathStr; got != want {
		t.Errorf("String() returned %q, want %q", got, want)
	}
}

func TestPath_Join(t *testing.T) {
	path := files.Path("/go/src/awesomeapp")

	if got, want := path.Join("cmd", "server"), files.Path("/go/src/awesomeapp/cmd/server"); got != want {
		t.Errorf("Join() returned %q, want %q", got, want)
	}
}
