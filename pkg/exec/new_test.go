package exec_test

import (
	"bytes"
	"fmt"

	"k8s.io/utils/exec"
)

func ExampleNew() {
	exec := exec.New()

	cmd := exec.Command("echo", "Bonjour!")
	buff := bytes.Buffer{}
	cmd.SetStdout(&buff)
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	fmt.Println(buff.String())
	// Output: Bonjour!
}
