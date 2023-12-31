package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"url-shortener/internal/config"
	mwLogger "url-shortener/internal/http-server/middleware/logger"
	"url-shortener/internal/lib/logger/cslog"
	"url-shortener/internal/lib/logger/handlers/slogpretty"
	"url-shortener/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	err, cfg := config.LoadCFG()
	if err != nil {
		log.Fatal(err)
	}

	log := setupLogger(cfg.Logger)

	log.Info("Starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("Debug message are enabled")

	storage, err := sqlite.New(cfg.Storage)

	if err != nil {
		log.Error(fmt.Sprintf("failed to init storage: '%s'", cfg.Storage), cslog.Err(err))
		os.Exit(1)
	}
	_ = storage

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// TODO: run server

}

func setupLogger(loggerCFG config.Logger) *slog.Logger {
	var log *slog.Logger
	switch loggerCFG.LogType {
	case "pretty":
		log = setupPrettySlog()
	case "text":
		log = slog.New(
			slog.NewTextHandler(
				os.Stdout, &slog.HandlerOptions{
					Level: slog.Level(loggerCFG.LogLevel),
				},
			),
		)
	case "json":
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout, &slog.HandlerOptions{
					Level: slog.Level(loggerCFG.LogLevel),
				},
			),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout, &slog.HandlerOptions{
					Level: slog.LevelInfo,
				},
			),
		)
	}
	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
