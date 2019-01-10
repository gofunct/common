package common

import "github.com/gofunct/common/ask"

type Ask interface {
	Q(q string) (string, error)
	TrueFalse(q string) (bool, error)
	YesNo(q string) (bool, error)
}

func NewAsk() Ask {
	return ask.API{}
}
