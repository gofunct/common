package log

import (
	"github.com/gofunct/common/utils"
	"go.uber.org/zap"
)

func (s *Service) FatalR(v ...interface{}) {
	s.Z.Fatal(utils.Fail+utils.Red(v), zap.Any("fatal", v))
}

func (s *Service) InfoB(v ...interface{}) {
	s.Z.Info(utils.Info+utils.Blue(v), zap.Any("info", v))
}

func (s *Service) DebugC(v ...interface{}) {
	s.Z.Info(utils.Debug+utils.Cyan(v), zap.Any("debug", v))
}

func (s *Service) SuccessG(v ...interface{}) {
	s.Z.Info(utils.Success+utils.Green(v), zap.Any("success", v))
}

func (s *Service) WarnY(v ...interface{}) {
	s.Z.Warn(utils.Warn+utils.Yello(v), zap.Any("warn", v))
}

func (s *Service) FireM(v ...interface{}) {
	s.Z.Warn(utils.Fire+utils.Magenta(v), zap.Any("fire", v))
}
