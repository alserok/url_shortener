package app

import (
	"context"
	"github.com/alserok/url_shortener/internal/cache"
	"github.com/alserok/url_shortener/internal/config"
	"github.com/alserok/url_shortener/internal/db"
	"github.com/alserok/url_shortener/internal/server"
	"github.com/alserok/url_shortener/internal/service"
	"github.com/alserok/url_shortener/pkg/logger"
	"os/signal"
	"syscall"
	"time"
)

func MustStart(cfg *config.Config) {
	log := logger.New(logger.Slog, cfg.Env)
	defer func() {
		log.Info("server has been stopped")
		_ = log.Close()
	}()

	c := cache.New(cache.Redis, cfg.Cache)
	defer func() {
		if err := c.Close(); err != nil {
			log.Warn("failed to close cache", logger.WithArg("error", err.Error()))
		}
	}()

	repo := db.New(cfg.DBType, cfg.DB)
	defer func() {
		if err := repo.Close(); err != nil {
			log.Warn("failed to close repo", logger.WithArg("error", err.Error()))
		}
	}()

	srvc := service.New(repo)

	srvr := server.New(cfg.ServerType, srvc, c, log)

	log.Info(
		"server is running",
		logger.WithArg("port", cfg.Port),
		logger.WithArg("env", cfg.Env),
		logger.WithArg("server_type", cfg.ServerType),
		logger.WithArg("db_type", cfg.DBType),
	)
	run(srvr, cfg.Port)
}

func run(s server.Server, port string) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go s.MustServe(port)

	<-ctx.Done()

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	s.Shutdown(ctx)
}
