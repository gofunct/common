package config

import (
	"fmt"
	"os"
	"testing"
)

func TestEnviron(t *testing.T) {
	s := os.Environ()
	fmt.Println(s)
	if len(s) == 0 {
		t.Fail()
	}
}

func TestPWD(t *testing.T) {
	wd, err := os.Getwd()
	fmt.Println(wd)
	fmt.Println(os.Getenv("PWD"))
	if err != nil {
		t.Fail()
	}
}

func TestBin(t *testing.T) {
	wd, err := os.Getwd()
	fmt.Println(wd)
	wd = os.Getenv("PWD")
	fmt.Println(wd)
	wd = os.Getenv("PWD") + "/bin"
	fmt.Println(wd)

	if err != nil {
		t.Fail()
	}
}
