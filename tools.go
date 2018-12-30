// +build tools

package tools

// tool dependencies
import (
	_ "github.com/haya14busa/reviewdog/cmd/reviewdog"
	_ "github.com/srvc/wraperr/cmd/wraperr"
	_ "golang.org/x/lint/golint"
	- "github.com/kisielk/errcheck"
	_ "github.com/mitchellh/gox"
	_ "honnef.co/go/tools"
	_ "mvdan.cc/unparam"
)
