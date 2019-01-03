package gocloud

import (
	"github.com/gofunct/common/logging"
	"github.com/gofunct/common/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"strings"
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
	listen 			string
	cloudSQLRegion    string
	runtimeConfigName string
}

var envFlag string

func BindGoCloudFlags(cmd *cobra.Command) error {
	var (
		config = viper.New()
		result string
		ask = ui.NewMessenger()
		err error
	)

	// Determine environment to set up based on flag.
	cf := new(cliFlags)
	config.SetConfigName("gocloud")
	config.AddConfigPath(".")
	config.AutomaticEnv()

	cmd.Flags().DurationVar(&cf.motdVarWaitTime, "runtime_var_wait", 5*time.Second, "polling frequency of message of the day")
	config.SetDefault("runtime_var_wait", 5*time.Second)
	logging.AddLoggingFlags(cmd)


	{
		result, err = ask.UI.Ask("what is environment do you want to run under?", " ")
		if err != nil {
			return err
		}
		envFlag = strings.ToLower(result)
		cmd.Flags().StringVar(&envFlag, "env", "local", "environment to run under")
		config.SetDefault("env", result)
	}
	{
		result, err = ask.UI.Ask("what port do you want to listen on?", " ")
		if err != nil {
			return err
		}
		cf.listen = strings.ToLower(result)
		cmd.Flags().String("listen", ":8080", "port to listen for HTTP on")
		config.SetDefault("listen", result)

	}
	{
		result, err = ask.UI.Ask("what bucket name do you want to setup?", " ")
		if err != nil {
			return err
		}
		cf.bucket = strings.ToLower(result)
		cmd.Flags().StringVar(&cf.bucket, "bucket", "", "bucket name")
		config.SetDefault("bucket", result)

	}
	{
		result, err = ask.UI.Ask("what is the database host or Cloud SQL instance name?", " ")
		if err != nil {
			return err
		}
		cf.dbHost = strings.ToLower(result)
		cmd.Flags().StringVar(&cf.dbHost, "db_host", "", "database host or Cloud SQL instance name")
		config.SetDefault("db_host", result)

	}
	{
		result, err = ask.UI.Ask("what is the database name?", " ")
		if err != nil {
			return err
		}
		cf.dbName = strings.ToLower(result)
		cmd.Flags().StringVar(&cf.dbName, "db_name", "", "database name")
		config.SetDefault("db_name", result)
	}
	{
		result, err = ask.UI.Ask("what is the database username?", " ")
		if err != nil {
			return err
		}
		cf.dbUser = strings.ToLower(result)
		cmd.Flags().StringVar(&cf.dbUser, "db_user", "guestbook", "database user")
		config.SetDefault("db_user", result)

	}
	{
		result, err = ask.UI.Ask("what is the database user password?", " ")
		if err != nil {
			return err
		}
		cf.dbPassword = strings.ToLower(result)
		cmd.Flags().StringVar(&cf.dbPassword, "db_password", "", "database user password")
		config.SetDefault("db_password", result)

	}
	{
		result, err = ask.UI.Ask("what is the runtime variable location?", " ")
		if err != nil {
			return err
		}
		cf.motdVar = strings.ToLower(result)
		cmd.Flags().StringVar(&cf.motdVar, "runtime_var", "", "runtime variable location")
		config.SetDefault("runtime_var", result)

	}
	{
		result, err = ask.UI.Ask("what region of the Cloud SQL instance (GCP only)?", " ")
		if err != nil {
			return err
		}
		cf.cloudSQLRegion = strings.ToLower(result)
		cmd.Flags().StringVar(&cf.cloudSQLRegion, "cloud_sql_region", "", "region of the Cloud SQL instance (GCP only)")
		config.SetDefault("cloud_sql_region", result)
	}
	{
		result, err = ask.UI.Ask("what is the runtime configurator config resource (GCP only)?", " ")
		if err != nil {
			return err
		}
		cf.runtimeConfigName = strings.ToLower(result)
		cmd.Flags().StringVar(&cf.runtimeConfigName, "runtime_config", "", "Runtime Configurator config resource (GCP only)")
		config.SetDefault("runtime_config", result)

	}

	config.BindPFlags(cmd.Flags())

	// If a config file is found, read it in.
	if err := config.ReadInConfig(); err != nil {
		log.Printf("%s, %s", "failed to read config file, writing defaults...", err)
		if err = config.WriteConfigAs("gocloud.yaml"); err != nil {
			return err
		}

	} else {
		log.Printf("%s, %s", "Using config file:", config.ConfigFileUsed())
		if err = config.WriteConfig(); err != nil {
			return err
		}
	}
	return nil
}
