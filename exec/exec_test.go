package exec

import (
	"github.com/gofunct/iio"
	"testing"
)

func TestExecutorNoArgs(t *testing.T) {
	ex := New("ls", iio.NewStdIO())

	out, err := ex.CombinedOutput()
	if err != nil {
		t.Errorf("expected success, got %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected no output, got %q", string(out))
	}
	out, err = ex.CombinedOutput()
	if err == nil {
		t.Errorf("expected failure, got nil error")
	}
	if len(out) != 0 {
		t.Errorf("expected no output, got %q", string(out))
	}
	ee, ok := err.(ExitError)
	if !ok {
		t.Errorf("expected an ExitError, got %+v", err)
	}
	if ee.Exited() {
		if code := ee.ExitStatus(); code != 1 {
			t.Errorf("expected exit status 1, got %d", code)
		}
	}
	out, err = ex.CombinedOutput()
	if err == nil {
		t.Errorf("expected failure, got nil error")
	}
	if ee, ok := err.(ExitError); ok {
		t.Errorf("expected non-ExitError, got %+v", ee)
	}
}

func TestExecutorWithArgs(t *testing.T) {
	ex := New("echo", iio.NewStdIO(), "hello world")

	out, err := ex.CombinedOutput()
	if err != nil {
		t.Errorf("expected success, got %+v", err)
	}
	if string(out) != "stdout\n" {
		t.Errorf("unexpected output: %q", string(out))
	}

	ex2 := New("/bin/sh", iio.NewStdIO(), "-c", "echo stderr > /dev/stderr")
	out, err = ex2.CombinedOutput()
	if err != nil {
		t.Errorf("expected success, got %+v", err)
	}
	if string(out) != "stderr\n" {
		t.Errorf("unexpected output: %q", string(out))
	}
}

func TestExecutableNotFound(t *testing.T) {
	cmd := New("fake_executable_name", iio.NewStdIO())
	_, err := cmd.CombinedOutput()
	if err != ErrExecutableNotFound {
		t.Errorf("cmd.CombinedOutput(): Expected error ErrExecutableNotFound but got %v", err)
	}

	{
		cmd := New("fake_executable_name", iio.NewStdIO())
		_, err := cmd.CombinedOutput()
		if err != ErrExecutableNotFound {
			t.Errorf("cmd.CombinedOutput(): Expected error ErrExecutableNotFound but got %v", err)
		}
	}
	{
		cmd := New("fake_executable_name", iio.NewStdIO())
		_, err = cmd.Output()
		if err != ErrExecutableNotFound {
			t.Errorf("cmd.Output(): Expected error ErrExecutableNotFound but got %v", err)
		}
	}

	cmd = New("fake_executable_name", iio.NewStdIO())
	err = cmd.Run()
	if err != ErrExecutableNotFound {
		t.Errorf("cmd.Run(): Expected error ErrExecutableNotFound but got %v", err)
	}
}
