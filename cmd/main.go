package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/OurLuv/geograkom/internal/config"
	"github.com/OurLuv/geograkom/internal/handler"
	"github.com/golang-cz/devslog"
	"gitlab.com/greyxor/slogor"
)

func main() {

	cfg := config.MustLoad()
	log := setupLogger(cfg.Env, cfg.LogLib)

	h := handler.NewHandler(log)
	r := h.InitRoutes()
	server := handler.NewServer(*cfg, r)

	log.Info("Starting application", slog.Any("cfg", cfg))

	if err := server.Start(); err != nil {
		log.Error("can't start server", slog.String("err", err.Error()))
	}

}

func setupLogger(env string, lib string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case "local":
		h := devHandlerLogger(lib)
		log = slog.New(h)
		slog.SetDefault(log)
	default:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelWarn,
		}))
	}
	return log
}

func devHandlerLogger(lib string) slog.Handler {
	var handler slog.Handler
	switch lib {
	case "devslog":
		slogOpts := &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}
		opts := &devslog.Options{
			HandlerOptions:    slogOpts,
			MaxSlicePrintSize: 4,
			SortKeys:          true,
			TimeFormat:        "[04:05]",
			NewLineAfterLog:   true,
			DebugColor:        devslog.Magenta,
		}
		handler = devslog.NewHandler(os.Stdout, opts)
	case "slogor":
		handler = slogor.NewHandler(os.Stderr, slogor.Options{
			TimeFormat: time.Stamp,
			Level:      slog.LevelDebug,
			ShowSource: true,
		})

	}
	return handler
}
