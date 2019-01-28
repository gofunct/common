package generator

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/gofunct/common/pkg/config"
	"github.com/gofunct/common/pkg/logger"
)

// Run generator
func Run(cfg *config.Config) {
	if cfg.Storage.Config.Name == "" {
		cfg.Storage.Config.Name = cfg.Name
	}
	if cfg.Storage.MySQL {
		cfg.Storage.Config.Driver = config.StorageMySQL
	}
	if cfg.Storage.Postgres {
		cfg.Storage.Config.Driver = config.StoragePostgres
	}
	logger.LogF("Base templates", copyTemplates(
		path.Join(cfg.Directories.Templates, config.Base),
		cfg.Directories.Service,
	))
	if cfg.API.Enabled {
		logger.LogF("Storage base templates", copyTemplates(
			path.Join(cfg.Directories.Templates, config.API, config.Base),
			cfg.Directories.Service,
		))
		if cfg.API.Gateway {
			logger.LogF("Gateway templates for API", copyTemplates(
				path.Join(cfg.Directories.Templates, config.API, config.APIGateway),
				cfg.Directories.Service,
			))
		}
	}
	if cfg.Storage.Enabled {
		logger.LogF("Storage base templates", copyTemplates(
			path.Join(cfg.Directories.Templates, config.Storage, config.Base),
			cfg.Directories.Service,
		))
		if cfg.Storage.Postgres {
			logger.LogF("Storage templates for postgres", copyTemplates(
				path.Join(cfg.Directories.Templates, config.Storage, config.StoragePostgres),
				cfg.Directories.Service,
			))
		}
		if cfg.Storage.MySQL {
			logger.LogF("Storage templates for mysql", copyTemplates(
				path.Join(cfg.Directories.Templates, config.Storage, config.StorageMySQL),
				cfg.Directories.Service,
			))
		}
	}
	if cfg.API.Enabled && cfg.Storage.Enabled && cfg.Contract {
		logger.LogF("Contract example templates", copyTemplates(
			path.Join(cfg.Directories.Templates, config.Contract, config.Base),
			cfg.Directories.Service,
		))
		if cfg.Storage.Postgres {
			logger.LogF("Contract templates for postgres", copyTemplates(
				path.Join(cfg.Directories.Templates, config.Contract, config.StoragePostgres),
				cfg.Directories.Service,
			))
		}
		if cfg.Storage.MySQL {
			logger.LogF("Contract templates for mysql", copyTemplates(
				path.Join(cfg.Directories.Templates, config.Contract, config.StorageMySQL),
				cfg.Directories.Service,
			))
		}
	}
	logger.LogF("Render templates", render(cfg))
	logger.LogF("Could not change directory", os.Chdir(cfg.Directories.Service))
	if cfg.API.Enabled && cfg.Storage.Enabled && cfg.Contract {
		log.Println("Prepare contracts:")
		logger.LogF("Generate contracts", Exec("make", "contracts"))
	}
	log.Println("Initialize vendors:")
	logger.LogF("Init dep", Exec("dep", "init", "-skip-tools"))
	logger.LogF("Tests", Exec("make", "check-all"))

	if cfg.GitInit {
		log.Println("Initialize Git repository:")
		logger.LogF("Init git", Exec("git", "init"))
		logger.LogF("Add repo files", Exec("git", "add", "--all"))
		logger.LogF("Initial commit", Exec("git", "commit", "-m", "'Initial commit'"))
	}
	fmt.Printf("New repository was created, use command 'cd %s'", cfg.Directories.Service)
}

// Exec runs the commands
func Exec(command ...string) error {
	execCmd := exec.Command(command[0], command[1:]...) // nolint: gosec
	execCmd.Stderr = os.Stderr
	execCmd.Stdout = os.Stdout
	execCmd.Stdin = os.Stdin
	return execCmd.Run()
}
