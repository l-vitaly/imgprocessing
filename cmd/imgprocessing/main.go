package main

import (
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/l-vitaly/imgprocessing/pkg/config"
	"github.com/l-vitaly/imgprocessing/pkg/log"
	"github.com/l-vitaly/imgprocessing/pkg/wire"
)

// Build vars.
var (
	GitHash   = ""
	BuildDate = ""
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() + 1)
}

func main() {
	logger := log.LoggerWrapDefault(
		kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stderr)),
		"imgproc",
		"IMG_PROC__LOG_LEVEL",
	)
	defer level.Info(logger).Log("msg", "goodbye")

	level.Info(logger).Log("version", GitHash, "builddate", BuildDate, "msg", "hello")

	cfg, err := config.Parse()
	if err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}

	err = os.MkdirAll(cfg.SavePath, 0755)
	if err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}

	h, cleanup, err := wire.SetupHTTPHandler(
		cfg.SavePath,
		logger,
	)
	if err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}
	defer cleanup()

	mux := http.NewServeMux()

	mux.Handle("/", h)
	http.Handle("/", accessControl(mux))

	errs := make(chan error)

	go func() {
		shutdown := make(chan os.Signal, 1)
		signal.Notify(
			shutdown,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGKILL,
			syscall.SIGTERM,
			syscall.SIGQUIT,
			syscall.SIGABRT,
			syscall.SIGSEGV,
		)

		<-shutdown

		errs <- nil
	}()

	go func() {
		level.Info(logger).Log("transport", "http", "address", cfg.BindAddr, "msg", "listening")
		errs <- http.ListenAndServe(cfg.BindAddr, nil)
	}()

	if err := <-errs; err != nil {
		level.Error(logger).Log("err", err)
	}

	level.Info(logger).Log("msg", "server terminated")
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}
