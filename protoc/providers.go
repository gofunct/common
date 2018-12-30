package protoc

import (
	"github.com/gofunct/common/bingen"
	"github.com/gofunct/common/bingen/tool"
	"github.com/gofunct/common/files"
	"github.com/gofunct/common/io"
	"github.com/gofunct/common/logging"
	"sync"

	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"go.uber.org/zap"
	"k8s.io/utils/exec"
	)

var (
	bingenCfg   *bingen.Config
	bingenCfgMu sync.Mutex

	toolRepo   tool.Repository
	toolRepoMu sync.Mutex
)

func ProvideGexConfig(
	fs afero.Fs,
	execer exec.Interface,
	io io.IO,
	rootDir files.RootDir,
) *bingen.Config {
	bingenCfgMu.Lock()
	defer bingenCfgMu.Unlock()
	if bingenCfg == nil {
		bingenCfg = &bingen.Config{
			OutWriter:  io.Out(),
			ErrWriter:  io.Err(),
			InReader:   io.In(),
			FS:         fs,
			Execer:     execer,
			WorkingDir: rootDir.String(),
			Verbose:    logging.IsVerbose() || logging.IsDebug(),
			Logger:     zap.NewStdLog(zap.L()),
		}
	}
	return bingenCfg
}

func ProvideToolRepository(bingenCfg *bingen.Config) (tool.Repository, error) {
	toolRepoMu.Lock()
	defer toolRepoMu.Unlock()
	if toolRepo == nil {
		var err error
		toolRepo, err = bingenCfg.Create()
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return toolRepo, nil
}

// WrapperSet is a provider set that includes gex things and Wrapper instance.
var WrapperSet = wire.NewSet(
	ProvideGexConfig,
	ProvideToolRepository,
	NewWrapper,
)

