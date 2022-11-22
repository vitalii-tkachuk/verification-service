package bootstrap

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/vitalii-tkachuk/verification-service/internal/application/verification/command"
	"github.com/vitalii-tkachuk/verification-service/internal/application/verification/query"
	"github.com/vitalii-tkachuk/verification-service/internal/domain/verification/service"
	"github.com/vitalii-tkachuk/verification-service/internal/infrastructure"
	"github.com/vitalii-tkachuk/verification-service/internal/infrastructure/bus"
	"github.com/vitalii-tkachuk/verification-service/internal/infrastructure/config"
	"github.com/vitalii-tkachuk/verification-service/internal/infrastructure/persistence/postgres"
	"github.com/vitalii-tkachuk/verification-service/internal/infrastructure/server"
)

var (
	ErrCannotParseConfig       = errors.New("cannot parse config")
	ErrCannotConnectToDatabase = errors.New("cannot connect to database")
)

// Run parses config environment variables, opens database connection and setup DI service for Buses
func Run() error {
	var cfg config.Config

	if err := envconfig.Process("", &cfg); err != nil {
		return fmt.Errorf("%s: %w", ErrCannotParseConfig, err)
	}

	db, err := sql.Open("postgres", cfg.PostgresDatabaseDsn())
	if err != nil {
		return fmt.Errorf("%s: %w", ErrCannotConnectToDatabase, err)
	}

	inMemoryCommandBus := bus.NewInMemoryCommandBus()
	queryBus := bus.NewQueryBus()

	verificationRepository := postgres.NewVerificationRepository(db, cfg.DatabaseTimeout)

	createVerificationService := service.NewCreateVerificationService(verificationRepository)
	approveVerificationService := service.NewApproveVerificationService(verificationRepository)
	declineVerificationService := service.NewDeclineVerificationService(verificationRepository)

	createVerificationCommandHandler := command.NewCreateVerificationCommandHandler(createVerificationService)
	approveVerificationCommandHandler := command.NewApproveVerificationCommandHandler(approveVerificationService)
	declineVerificationCommandHandler := command.NewDeclineVerificationCommandHandler(declineVerificationService)

	getVerificationByUUIDQueryHandler := query.NewGetVerificationByUUIDQueryHandler(verificationRepository)

	inMemoryCommandBus.Register(command.CreateVerificationCommandType, createVerificationCommandHandler)
	inMemoryCommandBus.Register(command.ApproveVerificationCommandType, approveVerificationCommandHandler)
	inMemoryCommandBus.Register(command.DeclineVerificationCommandType, declineVerificationCommandHandler)

	queryBus.Register(query.GetVerificationByUUIDQueryType, getVerificationByUUIDQueryHandler)

	application := infrastructure.NewApplication(inMemoryCommandBus, queryBus, validator.New())

	ctx, srv := server.NewServer(context.Background(), cfg, application)

	return srv.Run(ctx)
}
