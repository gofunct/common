package common

import (
	"github.com/gofunct/common/ask"
	"github.com/tcnksm/go-input"
)

type Ask interface {
	Query(q string) (string, error)
	TrueFalse(q string) (bool, error)
	YesNo(q string) (bool, error)
}

func NewAsk() Ask {
	return ask.API{
		Q: input.DefaultUI(),
	}
}
