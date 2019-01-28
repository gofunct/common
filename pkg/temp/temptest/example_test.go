package temptest

import (
	"errors"
	"fmt"
	"io"

	"k8s.io/utils/temp"
)

func TestedCode(dir temp.Directory) error {
	f, err := dir.NewFile("filename")
	if err != nil {
		return err
	}
	_, err = io.WriteString(f, "Bonjour!")
	if err != nil {
		return err
	}
	return dir.Delete()
}

func Example() {
	dir := FakeDir{}

	err := TestedCode(&dir)
	if err != nil {
		panic(err)
	}

	if dir.Deleted == false {
		panic(errors.New("Directory should have been deleted"))
	}

	if dir.Files["filename"] == nil {
		panic(errors.New(`"filename" should have been created`))
	}

	fmt.Println(dir.Files["filename"].Buffer.String())
	// Output: Bonjour!
}
