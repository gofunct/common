//+build wireinject

package cloud

import (
	"context"
	"database/sql"
	"github.com/gofunct/common/config"
	"github.com/google/wire"
	"github.com/gorilla/mux"
	"go.opencensus.io/trace"
	"gocloud.dev/blob"
	"gocloud.dev/health"
	"gocloud.dev/health/sqlhealth"
	"gocloud.dev/server"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func Serve(cfg *config.Service) {

	ctx := context.Background()
	var app *application
	var cleanup func()
	var err error
	switch cfg.Deploy {
	case "gcp":
		app, cleanup, err = SetupGCP(ctx, cfg)
	case "local":
		if cfg.DbHost == "" {
			cfg.DbHost = "localhost"
		}
		if cfg.DbPassword == "" {
			cfg.DbPassword = "password"
		}
		app, cleanup, err = SetupLocal(ctx, cfg)
	default:
		log.Fatalf("unknown --deploy=%s", cfg.Deploy)
	}
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	// Set up URL routes.
	r := mux.NewRouter()
	r.HandleFunc("/blob/{key:.+}", app.ServeBlob)

	// Listen and serve HTTP.
	log.Printf("Running, connected to %q cloud", cfg.Deploy)
	log.Fatal(app.srv.ListenAndServe(cfg.Lis, r))
}

var applicationSet = wire.NewSet(
	newApplication,
	appHealthChecks,
	trace.AlwaysSample,
)

type application struct {
	srv    *server.Server
	db     *sql.DB
	bucket *blob.Bucket
}

// newApplication creates a new application struct based on the backends and the message
// of the day variable.
func newApplication(srv *server.Server, db *sql.DB, bucket *blob.Bucket) *application {
	app := &application{
		srv:    srv,
		db:     db,
		bucket: bucket,
	}
	return app
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

func appHealthChecks(db *sql.DB) ([]health.Checker, func()) {
	dbCheck := sqlhealth.New(db)
	list := []health.Checker{dbCheck}
	return list, func() {
		dbCheck.Stop()
	}
}
