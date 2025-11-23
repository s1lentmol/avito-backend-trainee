package app

import (
	"fmt"
	"log/slog"

	"github.com/go-faster/errors"
	"github.com/silentmol/avito-backend-trainee/config"
	"github.com/silentmol/avito-backend-trainee/internal/controller/http"
	prrepo "github.com/silentmol/avito-backend-trainee/internal/pr/adapter/postgres"
	prusecase "github.com/silentmol/avito-backend-trainee/internal/pr/usecase"
	"github.com/silentmol/avito-backend-trainee/internal/storage"
	teamrepo "github.com/silentmol/avito-backend-trainee/internal/team/adapter/postgres"
	teamusecase "github.com/silentmol/avito-backend-trainee/internal/team/usecase"
	userrepo "github.com/silentmol/avito-backend-trainee/internal/user/adapter/postgres"
	userusecase "github.com/silentmol/avito-backend-trainee/internal/user/usecase"
	"github.com/silentmol/avito-backend-trainee/migrator"
)

func Run() error {
	slog.Info("starting application")

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", slog.Any("error", err))
		return errors.Wrap(err, "config")
	}
	slog.Info("config loaded",
		slog.String("app_name", cfg.App.Name),
		slog.String("app_port", cfg.App.Port),
		slog.String("db_host", cfg.DB.Host),
		slog.String("db_port", cfg.DB.Port),
		slog.String("db_name", cfg.DB.Name),
	)

	conn, err := storage.GetConnect(cfg.GetDSN())
	if err != nil {
		slog.Error("failed to connect to database", slog.Any("error", err))
		return errors.Wrap(err, "db conn")
	}
	defer conn.Close()
	slog.Info("database connection established")

	if err := migrator.Migrate(cfg.GetDSN()); err != nil {
		slog.Error("failed to apply migrations", slog.Any("error", err))
		return errors.Wrap(err, "db migrate")
	}
	slog.Info("database migrations applied")

	userRepo := userrepo.NewUserRepository(conn)
	teamRepo := teamrepo.NewTeamRepository(conn)
	prRepo := prrepo.NewPRRepository(conn)

	userUsecase := userusecase.NewUserUsecase(userRepo)
	teamUsecase := teamusecase.NewTeamUsecase(teamRepo)
	prUsecase := prusecase.NewPRUsecase(prRepo, userRepo, teamRepo)

	handle := http.NewHandler(userUsecase, teamUsecase, prUsecase)

	app := getRouter(handle, cfg.App.Name)
	slog.Info("starting http server", slog.String("port", cfg.App.Port))
	if err := app.Listen(":" + cfg.App.Port); err != nil {
		slog.Error("http server listen error", slog.Any("error", err))
		return fmt.Errorf("listen port: %w", err)
	}

	slog.Info("application stopped")
	return nil
}
