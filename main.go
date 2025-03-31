package main

import (
	"net/http"
	"os"
	"time"

	"github.com/Luzifer/doc-render/pkg/api"
	"github.com/Luzifer/doc-render/pkg/frontend"
	"github.com/Luzifer/doc-render/pkg/persist/mem"
	"github.com/Luzifer/doc-render/pkg/persist/redis"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	httpHelper "github.com/Luzifer/go_helpers/v2/http"
	"github.com/Luzifer/rconfig/v2"
)

var (
	cfg = struct {
		Listen          string `flag:"listen" default:":3000" description:"Port/IP to listen on"`
		LogLevel        string `flag:"log-level" default:"info" description:"Log level (debug, info, warn, error, fatal)"`
		PersistTo       string `flag:"persist-to" default:"disable" description:"Where to store server-side templates (disable, redis)"`
		SourceSetFolder string `flag:"source-set-folder" default:"source" description:"Where to find the templates to render"`
		TexAPIJobURL    string `flag:"tex-api-job-url" default:"" description:"Where to find the job endpoint of the TeX-API"`
		VersionAndExit  bool   `flag:"version" default:"false" description:"Prints current version and exits"`
	}{}

	version = "dev"
)

func initApp() error {
	rconfig.AutoEnv(true)
	if err := rconfig.ParseAndValidate(&cfg); err != nil {
		return errors.Wrap(err, "parsing cli options")
	}

	l, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		return errors.Wrap(err, "parsing log-level")
	}
	logrus.SetLevel(l)

	return nil
}

func main() {
	var err error
	if err = initApp(); err != nil {
		logrus.WithError(err).Fatal("initializing app")
	}

	if cfg.VersionAndExit {
		logrus.WithField("version", version).Info("doc-render")
		os.Exit(0)
	}

	r := mux.NewRouter()

	apiOpts := []api.Option{
		api.WithSourceSetDir(cfg.SourceSetFolder),
		api.WithTexAPIJobURL(cfg.TexAPIJobURL),
	}

	switch cfg.PersistTo {
	case "disable", "":
		// Nothing to do, persistence is disabled

	case "mem":
		apiOpts = append(apiOpts, api.WithPersistBackend(mem.New()))

	case "redis":
		backend, err := redis.New()
		if err != nil {
			logrus.WithError(err).Fatal("creating redis backend")
		}
		apiOpts = append(apiOpts, api.WithPersistBackend(backend))

	default:
		logrus.Fatal("invalid persist-to backend")
	}

	apiHandler := api.New(apiOpts...)
	apiHandler.Register(r)

	frontendHandler := frontend.New()
	frontendHandler.Register(r)

	var hdl http.Handler = r
	hdl = httpHelper.NewHTTPLogHandlerWithLogger(hdl, logrus.StandardLogger())

	srv := http.Server{
		Addr:              cfg.Listen,
		Handler:           hdl,
		ReadHeaderTimeout: time.Second,
	}

	logrus.WithField("version", version).WithField("addr", cfg.Listen).Info("doc-render started")
	if err = srv.ListenAndServe(); err != nil {
		logrus.WithError(err).Fatal("listening for HTTP traffic")
	}
}
