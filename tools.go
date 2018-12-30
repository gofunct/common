// +build tools

package main

// tool dependencies
import (
	_ "github.com/haya14busa/reviewdog/cmd/reviewdog"
	_ "github.com/srvc/wraperr/cmd/wraperr"
	_ "golang.org/x/lint/golint"
)
