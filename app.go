package common

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gofunct/common/ask"
	"github.com/gofunct/common/fs"
	"github.com/gofunct/common/hack"
	"github.com/gofunct/common/log"
	"github.com/gofunct/common/render"
	"github.com/gofunct/iio"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gocloud.dev/blob"
	"gocloud.dev/server"
	"gopkg.in/pipe.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"time"
)

type application struct {
	srv      *server.Server
	db       *sql.DB
	bucket   *blob.Bucket
	Config   *Config
	FS       *fs.Service
	Scripter *hack.Service
	Q        *ask.Service
	Renderer *render.Service
	L        *log.Service
	IO       *iio.Service
	Router   *mux.Router
}

func (a *application) SetupLocalDb() error {
	image := "mysql:5.6"

	zap.L().Debug("Starting container running MySQL")
	dockerArgs := []string{"run", "--rm"}
	if a.Config.Container != "" {
		dockerArgs = append(dockerArgs, "--name", a.Config.Container)
	}
	dockerArgs = append(dockerArgs,
		"--env", "MYSQL_DATABASE="+a.Config.DbName,
		"--env", "MYSQL_ROOT_PASSWORD="+a.Config.DbPassword,
		"--detach",
		"--publish", "3306:3306",
		image)
	cmd := exec.Command("docker", dockerArgs...)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("running %v: %v: %s", cmd.Args, err, out)
	}
	containerID := strings.TrimSpace(string(out))
	defer func() {
		zap.L().Debug("killing", zap.String("container", containerID))
		stop := exec.Command("docker", "kill", containerID)
		stop.Stderr = os.Stderr
		if err := stop.Run(); err != nil {
			zap.L().Debug("failed to kill db container", zap.Error(err))

		}
	}()

	// Stop the container on Ctrl-C.
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		// TODO(ijt): Handle SIGTERM.
		signal.Notify(c, os.Interrupt)
		<-c
		cancel()
	}()

	nap := 10 * time.Second
	zap.L().Debug("Waiting %v for database to come up", zap.Duration("nap", nap))
	select {
	case <-time.After(nap):
		// ok
	case <-ctx.Done():
		return errors.New("interrupted while napping")
	}
	zap.L().Debug("Initializing database schema and users")
	schema, err := ioutil.ReadFile(filepath.Join(a.Config.Source, "schema.sql"))
	if err != nil {
		return fmt.Errorf("reading schema: %v", err)
	}
	roles, err := ioutil.ReadFile(filepath.Join(a.Config.Source, "roles.sql"))
	if err != nil {
		return fmt.Errorf("reading roles: %v", err)
	}
	tooMany := 10
	var i int
	for i = 0; i < tooMany; i++ {
		mySQL := `mysql -h"${MYSQL_PORT_3306_TCP_ADDR?}" -P"${MYSQL_PORT_3306_TCP_PORT?}" -uroot -ppassword guestbook`
		p := pipe.Line(
			pipe.Read(strings.NewReader(string(schema)+string(roles))),
			pipe.Exec("docker", "run", "--rm", "--interactive", "--link", containerID+":mysql", image, "sh", "-c", mySQL),
		)
		if _, stderr, err := pipe.DividedOutput(p); err != nil {
			zap.L().Debug("Failed to seed database, retrying", zap.Any("stderr", stderr))
			select {
			case <-time.After(time.Second):
				continue
			case <-ctx.Done():
				return errors.New("interrupted while napping in between database seeding attempts")
			}
		}
		break
	}
	if i == tooMany {
		return fmt.Errorf("gave up after %d tries to seed database", i)
	}
	zap.L().Debug("Database running at localhost:3306")
	attach := exec.CommandContext(ctx, "docker", "attach", containerID)
	attach.Stdout = os.Stdout
	attach.Stderr = os.Stderr
	if err := attach.Run(); err != nil {
		return fmt.Errorf("running %v: %q", attach.Args, err)
	}

	return nil
}

func (app *application) RunLocalDb() {

	if err := app.SetupLocalDb(); err != nil {
		zap.L().Fatal("failed to run local db", zap.Error(errors.WithStack(err)))
	}
}
