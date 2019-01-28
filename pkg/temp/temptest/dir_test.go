package temptest

import (
	"io"
	"testing"
)

func TestFakeDir(t *testing.T) {
	d := &FakeDir{}

	f, err := d.NewFile("ONE")
	if err != nil {
		t.Fatal(err)
	}

	n, err := io.WriteString(f, "Bonjour!")
	if n != 8 || err != nil {
		t.Fatalf(
			`WriteString(f, "Bonjour!") = (%v, %v), expected (%v, %v)`,
			n, err,
			0, nil,
		)
	}
	if got := d.Files["ONE"].Buffer.String(); got != "Bonjour!" {
		t.Fatalf(`file content is %q, expected "Bonjour!"`, got)
	}

	f, err = d.NewFile("ONE")
	if err == nil {
		t.Fatal("Same file could be created twice.")
	}

	err = d.Delete()
	if err != nil {
		t.Fatal(err)
	}

	err = d.Delete()
	if err == nil {
		t.Fatal("FakeDir could be deleted twice.")
	}

	f, err = d.NewFile("TWO")
	if err == nil {
		t.Fatal("NewFile could be created in deleted dir")
	}

	if !d.Deleted {
		t.Fatal("FakeDir should be deleted.")
	}
}
