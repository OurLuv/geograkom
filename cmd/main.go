package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OurLuv/geograkom/internal/config"
	"github.com/OurLuv/geograkom/internal/handler"
	"github.com/OurLuv/geograkom/internal/service"
	"github.com/OurLuv/geograkom/internal/storage"
	"github.com/golang-cz/devslog"
	"gitlab.com/greyxor/slogor"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.MustLoad()
	log := setupLogger(cfg.Env, cfg.LogLib)

	conn, err := storage.NewConn(*cfg)
	if err != nil {
		log.Error("panic", slog.String("msg", err.Error()))
		panic(conn)
	}
	defer conn.Close()

	repo := storage.NewRouteStorage(conn, log)
	s := service.NewRouteServcie(repo)
	h := handler.NewHandler(s, log)
	r := h.InitRoutes()
	server := handler.NewServer(*cfg, r)

	log.Info("Starting application", slog.Any("cfg", cfg))

	go func() {
		if err := server.Start(); err != nil {
			log.Debug("server is off", slog.String("err", err.Error()))
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-ctx.Done()
	log.Info("Shutting down")
	err = server.ShutDown()
	if err != nil {
		log.Error("Error while shutting down the server", slog.String("err", err.Error()))
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
