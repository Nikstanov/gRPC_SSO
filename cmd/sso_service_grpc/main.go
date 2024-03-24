package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sso_service_grps/internal/app"
	"sso_service_grps/internal/config"
	"syscall"
)

const (
	envLocal = "local"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	log := setupLogger(cfg.Env)
	log.Info("starting app", slog.Any("config", cfg))

	application := app.New(log, cfg.GRPC.Port, cfg.DB, cfg.TokenTTL)
	go application.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	log.Info("stopping app", slog.String("signal", sign.String()))
	application.GRPCServer.Stop()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}
	return log
}
