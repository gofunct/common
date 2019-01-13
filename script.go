package common

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
	"os/exec"
	"os/user"
)

func (a *application) SetStdin(in io.Reader) {
	a.IO.InR = in
}

func (a *application) SetStdout(out io.Writer) {
	a.IO.OutW = out
}

func (a *application) SetStderr(out io.Writer) {
	a.IO.ErrW = out
}

func (a *application) RequireRoot() error {
	u, err := user.Current()
	if err != nil {
		return errors.Wrap(err, "failed to look up current user")
	}
	if u.Name != "root" {
		return errors.Wrap(err, "root user is required")
	}

	return nil
}
func (a *application) Gex(args ...string) []byte {
	out, err := exec.Command("gex", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run gex", zap.Strings("args", args), zap.Error(err))
	}
	return out
}

func (a *application) Git(args ...string) []byte {
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run git", zap.Strings("args", args), zap.Error(err))
	}
	return out
}

func (a *application) Gcloud(args ...string) []byte {
	out, err := exec.Command("gcloud", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run gcloud", zap.Strings("args", args), zap.Error(err))
	}
	return out
}

func (a *application) Bash(args ...string) []byte {
	out, err := exec.Command("bash", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run bash", zap.Strings("args", args), zap.Error(err))
	}
	return out
}

func (a *application) Go(args ...string) []byte {
	out, err := exec.Command("go", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run go", zap.Strings("args", args), zap.Error(err))
	}
	return out
}

func (a *application) Stencil(args ...string) []byte {
	out, err := exec.Command("stencil", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run stencil", zap.Strings("args", args), zap.Error(err))
	}
	return out
}

func (a *application) Make(args ...string) []byte {

	out, err := exec.Command("make", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run make", zap.Strings("make", args), zap.Error(err))
	}
	return out
}

func (a *application) Docker(args ...string) []byte {

	out, err := exec.Command("docker", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run docker", zap.Strings("args", args), zap.Error(err))
	}
	return out
}

func (a *application) Kubectl(args ...string) []byte {

	out, err := exec.Command("kubectl", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run kubectl", zap.Strings("args", args), zap.Error(err))
	}
	return out
}

func (a *application) Protoc(args ...string) []byte {

	out, err := exec.Command("protoc", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run protoc", zap.Strings("args", args), zap.Error(err))
	}
	return out
}

func (a *application) Aws(args ...string) []byte {

	out, err := exec.Command("aws", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run protoc", zap.Strings("args", args), zap.Error(err))
	}
	return out
}

func (a *application) Terraform(args ...string) []byte {

	out, err := exec.Command("terraform", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run terraform", zap.Strings("args", args), zap.Error(err))
	}
	return out
}

func (a *application) Ansible(args ...string) []byte {

	out, err := exec.Command("ansible", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run ansible", zap.Strings("args", args), zap.Error(err))
	}
	return out
}

func (a *application) Wire(args ...string) []byte {

	out, err := exec.Command("wire", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run wire", zap.Strings("args", args), zap.Error(err))
	}
	return out
}

func (a *application) Helm(args ...string) []byte {

	out, err := exec.Command("helm", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run helm", zap.Strings("args", args), zap.Error(err))
	}
	return out
}

func (a *application) Sed(args ...string) []byte {

	out, err := exec.Command("sed", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run sed", zap.Strings("args", args), zap.Error(err))
	}
	return out
}

func (a *application) Grep(args ...string) []byte {

	out, err := exec.Command("grep", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run grep", zap.Strings("args", args), zap.Error(err))
	}
	return out
}

func (a *application) Hero(args ...string) []byte {

	out, err := exec.Command("hero", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run hero", zap.Strings("args", args), zap.Error(err))
	}
	return out
}
