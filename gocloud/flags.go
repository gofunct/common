package gocloud

import (
	"github.com/spf13/cobra"
	"time"
)

type cliFlags struct {
	bucket          string
	dbHost          string
	dbName          string
	dbUser          string
	dbPassword      string
	motdVar         string
	motdVarWaitTime time.Duration

	cloudSQLRegion    string
	runtimeConfigName string
}

var envFlag string

func BindGoCloudFlags(cmd *cobra.Command) {
	// Determine environment to set up based on flag.
	cf := new(cliFlags)
	cmd.Flags().StringVar(&envFlag, "env", "local", "environment to run under")
	cmd.Flags().String("listen", ":8080", "port to listen for HTTP on")
	cmd.Flags().StringVar(&cf.bucket, "bucket", "", "bucket name")
	cmd.Flags().StringVar(&cf.dbHost, "db_host", "", "database host or Cloud SQL instance name")
	cmd.Flags().StringVar(&cf.dbName, "db_name", "guestbook", "database name")
	cmd.Flags().StringVar(&cf.dbUser, "db_user", "guestbook", "database user")
	cmd.Flags().StringVar(&cf.dbPassword, "db_password", "", "database user password")
	cmd.Flags().StringVar(&cf.motdVar, "motd_var", "", "message of the day variable location")
	cmd.Flags().DurationVar(&cf.motdVarWaitTime, "motd_var_wait_time", 5*time.Second, "polling frequency of message of the day")
	cmd.Flags().StringVar(&cf.cloudSQLRegion, "cloud_sql_region", "", "region of the Cloud SQL instance (GCP only)")
	cmd.Flags().StringVar(&cf.runtimeConfigName, "runtime_config", "", "Runtime Configurator config resource (GCP only)")
}
