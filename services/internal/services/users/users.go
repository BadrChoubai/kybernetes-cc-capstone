package users

import (
	"context"

	"github.com/badrchoubai/services/internal/config"
	"github.com/badrchoubai/services/internal/database"
	"github.com/badrchoubai/services/internal/observability/logging"
	"github.com/badrchoubai/services/internal/service"
)

func NewUsersService(ctx context.Context, cfg *config.AppConfig) (*service.Service, error) {
	logger, err := logging.NewLogger()
	if err != nil {
		return nil, err
	}

	db, err := database.NewDatabase(cfg)
	if err != nil {
		logger.Error("establishing database connection", err)
		return nil, err
	}

	svc := service.NewService(
		ctx,
		service.WithName("users-service"),
		service.WithURL("/api/v1/users"),
		service.WithLogger(logger),
		service.WithDatabase(db),
	)

	addRoutes(svc)
	return svc, nil
}
