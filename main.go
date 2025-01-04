package main

import (
	"net/http"
	"os"
	"time"

	"github.com/Luzifer/doc-render/pkg/api"
	"github.com/Luzifer/doc-render/pkg/frontend"
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

	apiHandler := api.New(
		api.WithSourceSetDir(cfg.SourceSetFolder),
		api.WithTexAPIJobURL(cfg.TexAPIJobURL),
	)
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
