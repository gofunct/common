package files

import (
	"github.com/gofunct/common/errors"
	"github.com/serenize/snaker"
	"path/filepath"
	"strings"
)

type ProtoParams struct {
	Proto struct {
		Path    string
		Package string
	}
	PbGo struct {
		Package    string
		ImportName string
	}
}

func BuildProtoParams(path string, rootDir RootDir, protoOutDir string, pkg string) (out ProtoParams, err error) {
	if protoOutDir == "" {
		err = errors.New("protoOutDir is required")
		return
	}

	// github.com/foo/bar
	importPath, err := GetImportPath(rootDir.String())
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	// path => baz/qux/quux
	path = strings.Replace(snaker.CamelToSnake(path), "-", "_", -1)

	// baz/qux
	packagePath := filepath.Dir(path)

	// api/baz/qux
	pbgoPackagePath := filepath.Join(protoOutDir, packagePath)
	// qux_pb
	pbgoPackageName := filepath.Base(pbgoPackagePath) + "_pb"

	if packagePath == "." {
		pbgoPackagePath = protoOutDir
		pbgoPackageName = filepath.Base(pbgoPackagePath) + "_pb"
	}

	protoPackage := pkg
	if protoPackage == "" {
		protoPackageChunks := []string{}
		for _, pkg := range strings.Split(filepath.ToSlash(filepath.Join(importPath, protoOutDir)), "/") {
			chunks := strings.Split(pkg, ".")
			for i := len(chunks) - 1; i >= 0; i-- {
				protoPackageChunks = append(protoPackageChunks, chunks[i])
			}
		}
		// com.github.foo.bar.baz.qux
		protoPackage = strings.Join(protoPackageChunks, ".")
	}
	if dir := filepath.Dir(path); dir != "." {
		protoPackage = protoPackage + "." + strings.Replace(dir, string(filepath.Separator), ".", -1)
	}
	protoPackage = strings.Replace(protoPackage, "-", "_", -1)

	out.Proto.Path = path
	out.Proto.Package = strings.Replace(protoPackage, "-", "_", -1)
	out.PbGo.Package = filepath.ToSlash(filepath.Join(importPath, pbgoPackagePath))
	out.PbGo.ImportName = pbgoPackageName

	return
}


type Params struct {
	ProtoDir    string
	ProtoOutDir string
	ServerDir   string
	Path        string
	ServiceName string
	Methods     []MethodParams
	Proto       ProtoParams
	PbGo        PbGoParams
	Go          GoParams
}

type ProtoParameters struct {
	Package  string
	Imports  []string
	Messages []MethodMessage
}

type PbGoParams struct {
	PackageName string
	PackagePath string
}

type GoParams struct {
	Package     string
	Imports     []string
	TestImports []string
	ServerName  string
	StructName  string
}

type MethodsParams struct {
	Methods      []MethodParams
	ProtoImports []string
	GoImports    []string
	Messages     []MethodMessage
}

type MethodParams struct {
	Method         string
	HTTP           MethodHTTPParams
	requestCommon  string
	requestGo      string
	requestProto   string
	responseCommon string
	responseGo     string
	responseProto  string
}

func (p *MethodParams) RequestGo(pkg string) string {
	if p.requestGo == "" {
		return pkg + "." + p.requestCommon
	}
	return p.requestGo
}

func (p *MethodParams) RequestProto() string {
	if p.requestProto == "" {
		return p.requestCommon
	}
	return p.requestProto
}

func (p *MethodParams) ResponseGo(pkg string) string {
	if p.responseGo == "" {
		return pkg + "." + p.responseCommon
	}
	return p.responseGo
}

func (p *MethodParams) ResponseProto() string {
	if p.responseProto == "" {
		return p.responseCommon
	}
	return p.responseProto
}

type MethodMessage struct {
	Name   string
	Fields []MethodMessageField
}

type MethodMessageField struct {
	Name     string
	Type     string
	Repeated bool
	Tag      uint
}

type MethodHTTPParams struct {
	Method string
	Path   string
	Body   string
}
