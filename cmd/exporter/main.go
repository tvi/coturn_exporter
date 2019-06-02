package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/tvi/coturn_exporter/coturn"

	flags "github.com/jessevdk/go-flags"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var opts struct {
	LogLevel string `long:"log-level" env:"LOG_LEVEL" description:"Log level" default:"info"`
	BindAddr string `long:"bind-address" env:"BIND_ADDRESS" default:":9596" description:"address for binding"`

	Username string `long:"username" short:"u" env:"USERNAME" description:"admin username for coturn server" required:"true"`
	Password string `long:"password" short:"p" env:"PASSWORD" description:"admin password for coturn server" required:"true"`
	Endpoint string `long:"endpoint" short:"e" env:"ENDPOINT" description:"admin interface endpoint" required:"true"`
}

func main() {
	parser := flags.NewParser(&opts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		// If the error was from the parser, then we can simply return
		// as Parse() prints the error already
		if _, ok := err.(*flags.Error); ok {
			os.Exit(1)
		}
		logrus.Fatalf("Error parsing flags: %v", err)
	}

	// Use log level
	level, err := logrus.ParseLevel(opts.LogLevel)
	if err != nil {
		logrus.Fatalf("Unknown log level %s: %v", opts.LogLevel, err)
	}
	logrus.SetLevel(level)

	// Set the log format to have a reasonable timestamp
	formatter := &logrus.TextFormatter{
		FullTimestamp: true,
	}
	logrus.SetFormatter(formatter)

	f, err := coturn.NewFetcher(opts.Username, opts.Password, opts.Endpoint)
	if err != nil {
		logrus.Fatalf("Could not initialize fetcher: %v", err)
	}

	if _, _, err = f.Fetch(context.Background()); err != nil {
		logrus.Fatalf("Could not do initial fetch: %v", err)
	}
	store, err := coturn.NewSessionStore(f)
	if err != nil {
		logrus.Fatalf("Could not initialize session store: %v", err)
	}

	go func(s *coturn.SessionStore) {
		for {
			if err := s.ReloadSessions(context.Background()); err != nil {
				logrus.Warn("could not reload sessions: %v", err)
			}
			time.Sleep(1 * time.Second)
		}
	}(store)

	ok := func(w http.ResponseWriter, _ *http.Request) { io.WriteString(w, "OK\n") }
	http.HandleFunc("/", ok)

	health := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		if _, _, err := f.Fetch(ctx); err != nil {
			logrus.Warnf("Got fetch error: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		io.WriteString(w, "OK\n")
	}
	http.HandleFunc("/health", health)

	http.Handle("/metrics", promhttp.Handler())

	logrus.Infof("Listening on %v\n", opts.BindAddr)
	if err := http.ListenAndServe(opts.BindAddr, nil); err != nil {
		logrus.Fatalf("Error binding: %v", err)
	}
}
