package authService

import (
	"context"
	"fmt"
	"log/slog"
	"sso_service_grps/internal/domain/models"
	"sso_service_grps/internal/utills"
	"time"
)

type AuthService struct {
	userProvider UserProvider
	userSaver    UserSaver
	appProvider  AppProvider
	log          *slog.Logger
	tokenTTL     time.Duration
}

func New(userProvider UserProvider, userSaver UserSaver, appProvider AppProvider, log *slog.Logger, tokenTTL time.Duration) *AuthService {
	return &AuthService{userProvider: userProvider, userSaver: userSaver, appProvider: appProvider, log: log, tokenTTL: tokenTTL}
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, password string) (uid int64, err error)
}
type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

func (a *AuthService) Register(ctx context.Context, email string, password string) (int64, error) {
	const op = "authService.Login"
	log := a.log.With(
		slog.String("op", op),
	)

	hashPassword, err := utills.HashPassword(password)
	if err != nil {
		log.Error("Failed to generate hash", slog.String("message", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	userID, err := a.userSaver.SaveUser(ctx, email, hashPassword)
	if err != nil {
		log.Error("Failed to save user", slog.String("message", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("User registered")
	return userID, nil
}

func (a *AuthService) Login(ctx context.Context, email string, password string, appID int32) (string, error) {
	const op = "authService.Register"
	log := a.log.With(
		slog.String("op", op),
	)

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		log.Error("Failed to find user", slog.String("message", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	res := utills.CheckPasswords(password, user.PassHash)
	if !res {
		log.Error("Failed to logIn user", slog.String("message", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	app, err := a.appProvider.App(ctx, int(appID))
	if err != nil {
		return "", err
	}

	token, err := utills.GenerateToken(user, app, a.tokenTTL)
	if err != nil {
		return "", err
	}
	log.Info("User logged")
	return token, nil
}
