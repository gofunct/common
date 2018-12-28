package logging

import "testing"

func TestEmojies(t *testing.T) {
	logger := NewLogger()

	logger.UI.Output("output")
	logger.UI.Running("running")
	logger.UI.Success("success")
	logger.UI.Log("log")
	logger.UI.Error("error")
	logger.UI.Info("info")
	logger.UI.Ask("?", "ask")
	logger.UI.Warn("warn")
}
