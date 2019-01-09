package ask

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/tcnksm/go-input"
	"go.uber.org/zap"
)

type Ask struct {
	*input.UI
}

func validateTF(q string) input.ValidateFunc {
	return func(ans string) error {
		zap.L().Debug("received response", zap.String("question", q), zap.String("answer", ans))
		if ans != "true" && ans != "false" {
			return fmt.Errorf("input must be true or false")
		}
		return nil
	}
}

func validateYN(q string) input.ValidateFunc {
	return func(ans string) error {
		zap.L().Debug("received response", zap.String("question", q), zap.String("answer", ans))
		if ans != "y" && ans != "n" {
			return fmt.Errorf("input must be y or n")
		}
		return nil
	}
}

func (u *Ask) YesNoBool(q string) (bool, error) {
	ans, err := u.Ask(fmt.Sprintf("%s [y/n]", q), &input.Options{
		HideOrder:    true,
		Loop:         true,
		ValidateFunc: validateYN(q),
	})
	if err != nil {
		return false, errors.WithStack(err)
	}
	return ans == "y", nil
}

func (u *Ask) TrueFalseBool(q string) (bool, error) {
	ans, err := u.Ask(fmt.Sprintf("%s [true/fase]", q), &input.Options{
		HideOrder:    true,
		Loop:         true,
		ValidateFunc: validateTF(q),
	})
	if err != nil {
		return false, errors.WithStack(err)
	}
	return ans == "true", nil
}

func (u *Ask) String(q string) (string, error) {
	ans, err := u.Ask(fmt.Sprintf("%s [text]", q), &input.Options{
		HideOrder: true,
		Loop:      true,
	})
	if err != nil {
		return "", errors.Wrapf(err, "failed to ask for input")
	}
	return ans, nil
}
