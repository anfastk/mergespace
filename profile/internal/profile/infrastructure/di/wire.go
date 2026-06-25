package di

import (
	inboundGrpc "github.com/anfastk/mergespace/profile/internal/profile/adapter/inbound/grpc/handler"

	inboundKafka "github.com/anfastk/mergespace/profile/internal/profile/adapter/inbound/kafka"

	"github.com/anfastk/mergespace/profile/internal/profile/adapter/outbound/postgres"

	"github.com/anfastk/mergespace/profile/internal/profile/application/usecase"

	"github.com/anfastk/mergespace/profile/internal/profile/infrastructure/config"
	"github.com/anfastk/mergespace/profile/internal/profile/infrastructure/database"
)

type App struct {
	ConsumerHandler *inboundKafka.ConsumerHandler

	GrpcHandler *inboundGrpc.ProfileHandler

	Config *config.Config
}

func BuildApp() *App {

	cfg := config.Load()

	db, err := database.NewPostgres(
		cfg.DatabaseURL,
	)
	if err != nil {
		panic(err)
	}

	repo := postgres.NewRepository(
		db,
	)

	profileUseCase := usecase.NewProfileUseCase(
		repo,
	)

	kafkaHandler := inboundKafka.NewConsumerHandler(
		profileUseCase,
	)

	grpcHandler := inboundGrpc.NewProfileHandler(
		profileUseCase,
	)

	return &App{
		ConsumerHandler: kafkaHandler,
		GrpcHandler:     grpcHandler,
		Config:          cfg,
	}
}
