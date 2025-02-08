package app

import (
	"context"
	"github.com/alserok/url_shortener/internal/config"
	"github.com/alserok/url_shortener/internal/db"
	"github.com/alserok/url_shortener/internal/server"
	"github.com/alserok/url_shortener/internal/service"
	"github.com/alserok/url_shortener/pkg/logger"
	"os/signal"
	"syscall"
)

func MustStart(cfg *config.Config) {
	log := logger.New(logger.Slog, cfg.Env)
	defer log.Info("server has been stopped")

	repo := db.New(cfg.DBType, cfg.DB)
	defer func() {
		_ = repo.Close()
	}()
	srvc := service.New(repo)
	srvr := server.New(cfg.ServerType, srvc, log)

	log.Info(
		"server is running",
		logger.WithArg("port", cfg.Port),
		logger.WithArg("env", cfg.Env),
		logger.WithArg("server type", cfg.ServerType),
		logger.WithArg("db type", cfg.DBType),
	)
	run(srvr, cfg.Port)
}

func run(s server.Server, port string) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go s.MustServe(port)

	<-ctx.Done()
}
