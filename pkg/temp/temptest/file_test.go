package temptest

import (
	"io"
	"testing"
)

func TestFakeFile(t *testing.T) {
	f := &FakeFile{}

	n, err := io.WriteString(f, "Bonjour!")
	if n != 8 || err != nil {
		t.Fatalf(
			`WriteString(f, "Bonjour!") = (%v, %v), expected (%v, %v)`,
			n, err,
			8, nil,
		)
	}

	err = f.Close()
	if err != nil {
		t.Fatal(err)
	}

	// File can't be closed twice.
	err = f.Close()
	if err == nil {
		t.Fatal("FakeFile could be closed twice")
	}

	// File is not writable after close.
	n, err = io.WriteString(f, "Bonjour!")
	if n != 0 || err == nil {
		t.Fatalf(
			`WriteString(f, "Bonjour!") = (%v, %v), expected (%v, %v)`,
			n, err,
			0, "non-nil",
		)
	}
}
