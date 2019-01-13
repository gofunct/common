//+build wireinject

package cloud

import (
	"context"
	"database/sql"
	"flag"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/wire"
	"github.com/gorilla/mux"
	"go.opencensus.io/trace"
	"gocloud.dev/blob"
	"gocloud.dev/health"
	"gocloud.dev/health/sqlhealth"
	"gocloud.dev/runtimevar"
	"gocloud.dev/server"
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

func Setup() {
	// Determine environment to set up based on flag.
	cf := new(cliFlags)
	flag.StringVar(&envFlag, "env", "local", "environment to run under")
	addr := flag.String("listen", ":8080", "port to listen for HTTP on")
	flag.StringVar(&cf.bucket, "bucket", "", "bucket name")
	flag.StringVar(&cf.dbHost, "db_host", "", "database host or Cloud SQL instance name")
	flag.StringVar(&cf.dbName, "db_name", "guestbook", "database name")
	flag.StringVar(&cf.dbUser, "db_user", "guestbook", "database user")
	flag.StringVar(&cf.dbPassword, "db_password", "", "database user password")
	flag.StringVar(&cf.motdVar, "motd_var", "", "message of the day variable location")
	flag.DurationVar(&cf.motdVarWaitTime, "motd_var_wait_time", 5*time.Second, "polling frequency of message of the day")
	flag.StringVar(&cf.cloudSQLRegion, "cloud_sql_region", "", "region of the Cloud SQL instance (GCP only)")
	flag.StringVar(&cf.runtimeConfigName, "runtime_config", "", "Runtime Configurator config resource (GCP only)")
	flag.Parse()

	ctx := context.Background()
	var app *application
	var cleanup func()
	var err error
	switch envFlag {
	case "gcp":
		app, cleanup, err = SetupGCP(ctx, cf)
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
	r.HandleFunc("/blob/{key:.+}", app.ServeBlob)

	// Listen and serve HTTP.
	log.Printf("Running, connected to %q cloud", envFlag)
	log.Fatal(app.srv.ListenAndServe(*addr, r))
}

// applicationSet is the Wire provider set for the Guestbook application that
// does not depend on the underlying platform.
var ApplicationSet = wire.NewSet(
	NewApplication,
	AppHealthChecks,
	trace.AlwaysSample,
)

// application is the main server struct for Guestbook. It contains the state of
// the most recently read message of the day.
type application struct {
	srv    *server.Server
	db     *sql.DB
	bucket *blob.Bucket

	// The following fields are protected by mu:
	mu   sync.RWMutex
	motd string // message of the day
}

// newApplication creates a new application struct based on the backends and the message
// of the day variable.
func NewApplication(srv *server.Server, db *sql.DB, bucket *blob.Bucket) *application {
	app := &application{
		srv:    srv,
		db:     db,
		bucket: bucket,
	}
	return app
}

// watchMOTDVar listens for changes in v and updates the app's message of the
// day. It is run in a separate goroutine.
func (app *application) WatchVar(v *runtimevar.Variable) {
	ctx := context.Background()
	for {
		snap, err := v.Watch(ctx)
		if err != nil {
			log.Printf("watch runtime variable: %v", err)
			continue
		}
		log.Println("updated runtime variable to", snap.Value)
		app.mu.Lock()
		app.motd = snap.Value.(string)
		app.mu.Unlock()
	}
}

// serveBlob handles a request for a static asset by retrieving it from a bucket.
func (app *application) ServeBlob(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	blobRead, err := app.bucket.NewReader(r.Context(), key, nil)
	if err != nil {
		// TODO(light): Distinguish 404.
		// https://github.com/google/go-cloud/issues/2
		log.Println("serve blob:", err)
		http.Error(w, "blob read error", http.StatusInternalServerError)
		return
	}
	// TODO(light): Get content type from blob storage.
	// https://github.com/google/go-cloud/issues/9
	switch {
	case strings.HasSuffix(key, ".png"):
		w.Header().Set("Content-Type", "image/png")
	case strings.HasSuffix(key, ".jpg"):
		w.Header().Set("Content-Type", "image/jpeg")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	w.Header().Set("Content-Length", strconv.FormatInt(blobRead.Size(), 10))
	if _, err = io.Copy(w, blobRead); err != nil {
		log.Println("Copying blob:", err)
	}
}

// appHealthChecks returns a health check for the database. This will signal
// to Kubernetes or other orchestrators that the server should not receive
// traffic until the server is able to connect to its database.
func AppHealthChecks(db *sql.DB) ([]health.Checker, func()) {
	dbCheck := sqlhealth.New(db)
	list := []health.Checker{dbCheck}
	return list, func() {
		dbCheck.Stop()
	}
}
