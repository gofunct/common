package exec

import (
	"go.uber.org/zap"
	"os/exec"
)

func (s *Scripter) Gex(args ...string) []byte {
	out, err := exec.Command("gex", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run gex", zap.Strings("args", args), zap.Error(err))
	}
	return out
}

func (s *Scripter) Git(args ...string) []byte {
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run git", zap.Strings("args", args), zap.Error(err))
	}
	return out
}

func (s *Scripter) Gcloud(args ...string) []byte {
	out, err := exec.Command("gcloud", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run gcloud", zap.Strings("args", args), zap.Error(err))
	}
	return out
}

func (s *Scripter) Bash(args ...string) []byte {
	out, err := exec.Command("bash", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run bash", zap.Strings("args", args), zap.Error(err))
	}
	return out
}

func (s *Scripter) Go(args ...string) []byte {
	out, err := exec.Command("go", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run go", zap.Strings("args", args), zap.Error(err))
	}
	return out
}

func (s *Scripter) Stencil(args ...string) []byte {
	out, err := exec.Command("stencil", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run stencil", zap.Strings("args", args), zap.Error(err))
	}
	return out
}

func (s *Scripter) Make(args ...string) []byte {

	out, err := exec.Command("make", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run make", zap.Strings("make", args), zap.Error(err))
	}
	return out
}

func (s *Scripter) Docker(args ...string) []byte {

	out, err := exec.Command("docker", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run docker", zap.Strings("args", args), zap.Error(err))
	}
	return out
}

func (s *Scripter) Kubectl(args ...string) []byte {

	out, err := exec.Command("kubectl", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run kubectl", zap.Strings("args", args), zap.Error(err))
	}
	return out
}

func (s *Scripter) Protoc(args ...string) []byte {

	out, err := exec.Command("protoc", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run protoc", zap.Strings("args", args), zap.Error(err))
	}
	return out
}

func (s *Scripter) Aws(args ...string) []byte {

	out, err := exec.Command("aws", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run protoc", zap.Strings("args", args), zap.Error(err))
	}
	return out
}

func (s *Scripter) Terraform(args ...string) []byte {

	out, err := exec.Command("terraform", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run terraform", zap.Strings("args", args), zap.Error(err))
	}
	return out
}

func (s *Scripter) Ansible(args ...string) []byte {

	out, err := exec.Command("ansible", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run ansible", zap.Strings("args", args), zap.Error(err))
	}
	return out
}

func (s *Scripter) Wire(args ...string) []byte {

	out, err := exec.Command("wire", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run wire", zap.Strings("args", args), zap.Error(err))
	}
	return out
}

func (s *Scripter) Helm(args ...string) []byte {

	out, err := exec.Command("helm", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run helm", zap.Strings("args", args), zap.Error(err))
	}
	return out
}

func (s *Scripter) Sed(args ...string) []byte {

	out, err := exec.Command("sed", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run sed", zap.Strings("args", args), zap.Error(err))
	}
	return out
}

func (s *Scripter) Grep(args ...string) []byte {

	out, err := exec.Command("grep", args...).Output()
	if err != nil {
		zap.L().Fatal("failed to run grep", zap.Strings("args", args), zap.Error(err))
	}
	return out
}
