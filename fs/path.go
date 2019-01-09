package fs

import "path/filepath"

// Path represents a filepath in filesystem.
type Path string

func (p Path) String() string { return string(p) }

// Join joins path elements to the path.
func (p Path) Join(elem ...string) Path {
	return Path(filepath.Join(append([]string{p.String()}, elem...)...))
}

type RootDir string

func (d RootDir) String() string { return string(d) }

func (d RootDir) Join(elem ...string) string {
	return filepath.Join(append([]string{d.String()}, elem...)...)
}

func (d RootDir) BinDir() string {
	return d.Join("bin")
}

func (d RootDir) CmdDir() string {
	return d.Join("cmd")
}

func (d RootDir) VendorDir() string {
	return d.Join("vendor")
}

func (d RootDir) ProtoDir() string {
	return d.Join("proto")
}

func (d RootDir) ApiDir() string {
	return d.Join("api")
}

func (d RootDir) DocsDir() string {
	return d.Join("docs")
}

func (d RootDir) StaticDir() string {
	return d.Join("static")
}
