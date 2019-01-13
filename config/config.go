package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/user"
	"runtime"
)

type Service struct {
	Bucket         string `mapstructure:"bucket"`
	DbHost         string `mapstructure:"db_host"`
	DbName         string `mapstructure:"db_name"`
	DbUser         string `mapstructure:"db_user"`
	DbPassword     string `mapstructure:"db_password"`
	CloudSqlRegion string `mapstructure:"cloud_sql_region"`
	Deploy         string `mapstructure:"deploy"`
	Lis            string `mapstructure:"lis"`
	*viper.Viper
}

func (s *Service) Bind(cmd *cobra.Command) {
	if err := s.Unmarshal(s); err != nil {
		log.Fatal("failed to unmarshal config", errors.WithStack(err))
	}

	cmd.Flags().StringVar(&s.Deploy, "deploy", "local", "environment to deploy to")
	cmd.Flags().StringVar(&s.Lis, "listen", ":8080", "port to listen for HTTP on")
	cmd.Flags().StringVar(&s.Bucket, "bucket", "", "bucket name")
	cmd.Flags().StringVar(&s.DbHost, "db_host", "", "database host or Cloud SQL instance name")
	cmd.Flags().StringVar(&s.DbName, "db_name", "", "database name")
	cmd.Flags().StringVar(&s.DbUser, "db_user", "", "database user")
	cmd.Flags().StringVar(&s.DbPassword, "db_password", "", "database user password")
	cmd.Flags().StringVar(&s.CloudSqlRegion, "cloud_sql_region", "", "region of the Cloud SQL instance (GCP only)")
	if err := s.BindPFlags(cmd.Flags()); err != nil {
		log.Fatal("failed to bind flags", errors.WithStack(err))
	}

	cmd.Annotations = s.Annotate()
	cmd.Version = cmd.Annotations["version"]
}

func (s *Service) Init() error {
	s.Viper = viper.GetViper()
	if s.Viper == nil {
		s.Viper = viper.New()
	}

	s.AutomaticEnv()

	ex, err := os.Executable()
	if err != nil {
		return errors.WithStack(err)
	}
	s.SetDefault("executable", ex)
	gr, err := os.Getgroups()

	if err != nil {
		return errors.WithStack(err)
	}
	s.SetDefault("groups", gr)
	host, err := os.Hostname()
	if err != nil {
		return errors.WithStack(err)
	}

	s.SetDefault("env", os.Environ())
	s.SetDefault("uid", os.Getuid())
	s.SetDefault("args", os.Args)
	s.SetDefault("host_name", host)
	s.SetDefault("pid", os.Getpid())
	s.SetDefault("goarch", runtime.GOARCH)
	s.SetDefault("compiler", runtime.Compiler)
	s.SetDefault("runtime_version", runtime.Version())
	s.SetDefault("goos", runtime.GOOS)
	s.SetDefault("goroot", runtime.GOROOT())
	usr, _ := user.Current()
	s.SetDefault("user", usr)
	return nil
}

func (s *Service) Annotate() map[string]string {
	settings := s.AllSettings()
	an := make(map[string]string)
	for k, v := range settings {
		if t, ok := v.(string); ok == true {
			an[k] = t
		}
	}
	return an
}
