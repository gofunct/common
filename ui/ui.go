package ui

import (
	"fmt"
	"github.com/gofunct/common/errors"
	iio "github.com/gofunct/common/io"
	"github.com/kyokomi/emoji"
	"gopkg.in/dixonwille/wlog.v2"
	"os"
	"strings"
	"sync"
)

// UI is an interface for intaracting with the terminal.
type UI interface {
	Section(msg string)
	Subsection(msg string)
	ItemSuccess(msg string)
	ItemSkipped(msg string)
	ItemFailure(msg string, errs ...error)
	Confirm(msg string) (bool, error)
}

var (
	ui   UI
	uiMu sync.Mutex
)

// UIInstance retuens a singleton UI instance.
func UIInstance(i iio.IO) UI {
	ui := wlog.New(i.In(), i.Out(), i.Err())

	pui := &wlog.PrefixUI{
		LogPrefix:     emoji.Sprint(":speech_balloon:"),
		OutputPrefix:  emoji.Sprint(":boom:"),
		SuccessPrefix: emoji.Sprint(":white_check_mark:"),
		InfoPrefix:    emoji.Sprint(":wave:"),
		ErrorPrefix:   emoji.Sprint(":x:"),
		WarnPrefix:    emoji.Sprint(":grimacing:"),
		RunningPrefix: emoji.Sprint(":fire:"),
		AskPrefix:     emoji.Sprint(":question:"),
		UI:            ui,
	}

	cui := wlog.AddConcurrent(pui)

	wlog.AddColor(wlog.Green, wlog.Red, wlog.BrightBlue, wlog.Blue, wlog.Yellow, wlog.BrightMagenta, wlog.Yellow, wlog.BrightGreen, wlog.BrightRed, cui)

	return &input{
		messenger: cui,
	}
}

// NewUI creates a new UI instance.
func NewUI(i iio.IO) UI {
	return &input{
		messenger: NewWLog(i),
		}
}

type input struct {
	inSection bool
	messenger wlog.UI
}

func (u *input) Section(msg string) {
	if u.inSection {
		u.inSection = false
	}
	u.Info("*SECTION* "+msg)
}

func (u *input) Subsection(msg string) {
	if u.inSection {
		u.inSection = false
	}
	u.Info("*SUBSECTION* "+ msg)
}

func (u *input) ItemSuccess(msg string) {
	u.inSection = true
	u.Success(msg)
	}

func (u *input) ItemSkipped(msg string) {
	u.inSection = true
	u.Info("*SKIPPED* "+ msg)
}

func (u *input) ItemFailure(msg string, errs ...error) {
	u.inSection = true
	for _, err := range errs {
		for _, s := range strings.Split(err.Error(), "\n") {
			u.Error(s)
		}
	}
}

func (u *input) Confirm(msg string) (bool, error) {
	ans := u.Ask(fmt.Sprintf("%s [y/n]", msg))

	if strings.Contains(ans, "y") {
		return true, nil
	} else if strings.Contains(ans, "n") {
		return false, nil
	}

	return false, errors.New("failed to confirm, answer must be y/n")
}

func (u input) Info(msg string) {
	u.messenger.Info(msg)
}

func (u input) Warn(msg string) {
	u.messenger.Warn(msg)
}

func (u input) Log(msg string) {
	u.messenger.Log(msg)
}

func (u input) Success(msg string) {
	u.messenger.Success(msg)
}

func (u input) Output(msg string) {
	u.messenger.Output(msg)
}

func (u input) Error(msg string) {
	u.messenger.Error(msg)
}

func (u input) Running(msg string) {
	u.messenger.Running(msg)
}

func (u input) Ask(msg string) string {
	s, err := u.messenger.Ask(msg, " ")
	if err != nil {
		u.Error(err.Error())
		os.Exit(1)
	}
	s = strings.ToLower(s)

	return s
}


func NewWLog(i iio.IO) *wlog.ConcurrentUI {
	ui := wlog.New(i.In(), i.Out(), i.Err())

	pui := &wlog.PrefixUI{
		LogPrefix:     emoji.Sprint(":speech_balloon:"),
		OutputPrefix:  emoji.Sprint(":boom:"),
		SuccessPrefix: emoji.Sprint(":white_check_mark:"),
		InfoPrefix:    emoji.Sprint(":wave:"),
		ErrorPrefix:   emoji.Sprint(":x:"),
		WarnPrefix:    emoji.Sprint(":grimacing:"),
		RunningPrefix: emoji.Sprint(":fire:"),
		AskPrefix:     emoji.Sprint(":question:"),
		UI:            ui,
	}

	cui := wlog.AddConcurrent(pui)

	wlog.AddColor(wlog.Green, wlog.Red, wlog.BrightBlue, wlog.Blue, wlog.Yellow, wlog.BrightMagenta, wlog.Yellow, wlog.BrightGreen, wlog.BrightRed, cui)

	return cui
}