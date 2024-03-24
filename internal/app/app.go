package app

import (
	"log/slog"
	grpcapp "sso_service_grps/internal/app/grpc"
	"sso_service_grps/internal/config"
	"sso_service_grps/internal/db"
	authService "sso_service_grps/internal/services/auth"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, conf config.DatabaseConfig, tokeTTL time.Duration) *App {
	storage, err := db.New(conf.User, conf.Password, conf.DBAddr, conf.Dbname, conf.Port)
	if err != nil {
		panic(err)
	}
	service := authService.New(storage, storage, storage, log, tokeTTL)
	grpcApp := grpcapp.New(log, service, grpcPort)

	return &App{grpcApp}
}
