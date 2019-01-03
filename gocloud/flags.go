package gocloud

import (
	"github.com/gofunct/common/ui"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"strings"
	"time"
	"context"
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

func NewGoCloudCommand(config *viper.Viper) (*cobra.Command, error) {
	var (
		cf = new(cliFlags)
		result string
		ask = ui.NewMessenger()
		err error
	)

	cmd := &cobra.Command{
		Use: "gocloud",
		Short: "cloud opts",
	}

	cmd.PersistentFlags().DurationVar(&cf.motdVarWaitTime, "runtime_var_wait", 5*time.Second, "polling frequency of message of the day")
	config.SetDefault("runtime_var_wait", 5*time.Second)


	{
		result, err = ask.UI.Ask("what is environment do you want to run under?(gcp or aws)", " ")
		if err != nil {
			return nil, err
		}
		envFlag = strings.ToLower(result)
		cmd.PersistentFlags().StringVar(&envFlag, "env", "local", "what is environment do you want to run under?(gcp or aws)")
		config.SetDefault("env", result)
	}
	{
		result, err = ask.UI.Ask("what port do you want to listen on?", " ")
		if err != nil {
			return nil, err
		}
		cmd.PersistentFlags().StringVar(&cf.listen, "listen", ":8080", "what port do you want to listen on?")
		config.SetDefault("listen", result)

	}
	{
		result, err = ask.UI.Ask("what bucket name do you want to setup?", " ")
		if err != nil {
			return nil, err
		}
		cf.bucket = strings.ToLower(result)
		cmd.PersistentFlags().StringVar(&cf.bucket, "bucket", "", "what bucket name do you want to setup?")
		config.SetDefault("bucket", result)

	}
	{
		result, err = ask.UI.Ask("what is the database host or Cloud SQL instance name?", " ")
		if err != nil {
			return nil, err
		}
		cf.dbHost = strings.ToLower(result)
		cmd.PersistentFlags().StringVar(&cf.dbHost, "db_host", "", "what is the database host or Cloud SQL instance name?")
		config.SetDefault("db_host", result)

	}
	{
		result, err = ask.UI.Ask("what is the database name?", " ")
		if err != nil {
			return nil, err
		}
		cf.dbName = strings.ToLower(result)
		cmd.PersistentFlags().StringVar(&cf.dbName, "db_name", "", "what is the database name?")
		config.SetDefault("db_name", result)
	}
	{
		result, err = ask.UI.Ask("what is the database username?", " ")
		if err != nil {
			return nil, err
		}
		cf.dbUser = strings.ToLower(result)
		cmd.PersistentFlags().StringVar(&cf.dbUser, "db_user", "guestbook", "what is the database username?")
		config.SetDefault("db_user", result)

	}
	{
		result, err = ask.UI.Ask("what is the database user password?", " ")
		if err != nil {
			return nil, err
		}
		cf.dbPassword = strings.ToLower(result)
		cmd.PersistentFlags().StringVar(&cf.dbPassword, "db_password", "", "what is the database user password?")
		config.SetDefault("db_password", result)

	}
	{
		result, err = ask.UI.Ask("what is the runtime variable location?", " ")
		if err != nil {
			return nil, err
		}
		cf.motdVar = strings.ToLower(result)
		cmd.PersistentFlags().StringVar(&cf.motdVar, "runtime_var", "", "what is the runtime variable location?")
		config.SetDefault("runtime_var", result)

	}
	{
		result, err = ask.UI.Ask("what region of the Cloud SQL instance (GCP only)?", " ")
		if err != nil {
			return nil, err
		}
		cf.cloudSQLRegion = strings.ToLower(result)
		cmd.PersistentFlags().StringVar(&cf.cloudSQLRegion, "cloud_sql_region", "", "region of the Cloud SQL instance (GCP only)")
		config.SetDefault("cloud_sql_region", result)
	}
	{
		result, err = ask.UI.Ask("what is the runtime configurator config resource (GCP only)?", " ")
		if err != nil {
			return nil, err
		}
		cf.runtimeConfigName = strings.ToLower(result)
		cmd.PersistentFlags().StringVar(&cf.runtimeConfigName, "runtime_config", "", "runtime Configurator config resource (GCP only)")
		config.SetDefault("runtime_config", result)

	}

	config.BindPFlags(cmd.PersistentFlags())

	init := &cobra.Command{
		Use: "init",
		Short: "initialize",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			var app *Application
			var cleanup func()
			var err error
			switch envFlag {
			case "gcp":
				app, cleanup, err = SetupGCP(ctx, cf)
			case "aws":
				if cf.dbPassword == "" {
					cf.dbPassword = "xyzzy"
				}
				app, cleanup, err = SetupAWS(ctx, cf)
			case "local":
				if cf.dbHost == "" {
					cf.dbHost = "localhost"
				}
				if cf.dbPassword == "" {
					cf.dbPassword = "xyzzy"
				}
				app, cleanup, err = SetupLocal(ctx, cf)
			default:
				log.Fatalf("unknown -env=%s", envFlag)
			}
			if err != nil {
				log.Fatal(err)
			}
			defer cleanup()

			// Set up URL routes.
			r := mux.NewRouter()
			//r.HandleFunc("/", app.index)
			//r.HandleFunc("/sign", app.sign)
			r.HandleFunc("/blob/{key:.+}", app.ServeBlob)

			// Listen and serve HTTP.
			log.Printf("Running, connected to %q cloud", envFlag)
			log.Fatal(app.srv.ListenAndServe(cf.listen, r))
		},
	}

	cmd.AddCommand(init)

	return cmd, nil
}
