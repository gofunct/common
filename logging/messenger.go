package logging

import (
	"github.com/kyokomi/emoji"
	"gopkg.in/dixonwille/wlog.v2"
	"os"
)

var (
	ui = wlog.New(os.Stdin, os.Stdout, os.Stderr)

	prefix = &wlog.PrefixUI{
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
)

type Messenger struct {
	UI wlog.UI
}

func NewMessenger() *Messenger {
	m := &Messenger{
		UI: prefix,
	}
	m.AddColor()

	return m
}

func (m *Messenger) AddColor() {
	m.UI = wlog.AddColor(wlog.Green, wlog.Red, wlog.BrightBlue, wlog.Blue, wlog.Yellow, wlog.BrightMagenta, wlog.Yellow, wlog.BrightGreen, wlog.BrightRed, m.UI)
}
