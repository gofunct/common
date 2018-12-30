package svcgen

import (
	"github.com/gofunct/common/files"
	"github.com/gofunct/common/proto/protoc"
	"github.com/gofunct/common/proto/service/params"
	"github.com/google/wire"
)

func ProvideParamsBuilder(rootDir files.RootDir, protocCfg *protoc.Config, grapiCfg *grapicmd.Config) params.Builder {
	return params.NewBuilder(
		rootDir,
		protocCfg.ProtosDir,
		protocCfg.OutDir,
		grapiCfg.Grapi.ServerDir,
		grapiCfg.Package,
	)
}

var Set = wire.NewSet(
	ProvideParamsBuilder,
	App{},
)
