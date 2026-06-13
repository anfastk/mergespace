package di

import (
	inboundKafka "github.com/anfastk/mergespace/profile/internal/profile/adapter/inbound/kafka"

	"github.com/anfastk/mergespace/profile/internal/profile/adapter/outbound/postgres"

	"github.com/anfastk/mergespace/profile/internal/profile/application/usecase"

	"github.com/anfastk/mergespace/profile/internal/profile/infrastructure/config"
	"github.com/anfastk/mergespace/profile/internal/profile/infrastructure/database"
)

type App struct {
	ConsumerHandler *inboundKafka.ConsumerHandler
	Config          *config.Config
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

	handler := inboundKafka.NewConsumerHandler(
		profileUseCase,
	)

	return &App{
		ConsumerHandler: handler,
		Config:          cfg,
	}

}
