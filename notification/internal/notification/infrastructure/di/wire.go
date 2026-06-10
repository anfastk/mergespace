package di

import (
	inboundKafka "github.com/anfastk/mergespace/notification/internal/notification/adapter/inbound/kafka"

	"github.com/anfastk/mergespace/notification/internal/notification/adapter/outbound/email"

	"github.com/anfastk/mergespace/notification/internal/notification/application/usecase"

	"github.com/anfastk/mergespace/notification/internal/notification/infrastructure/config"
)

type App struct {
	ConsumerHandler *inboundKafka.ConsumerHandler
	Config          *config.Config
}

func BuildApp() *App {

	cfg := config.Load()

	emailSender := email.NewSendGridSender(
		cfg.SendGrid.APIKey,
		cfg.SendGrid.From,
	)

	notificationUseCase :=
		usecase.NewNotificationUseCase(
			emailSender,
		)

	handler :=
		inboundKafka.NewConsumerHandler(
			notificationUseCase,
		)

	return &App{
		ConsumerHandler: handler,
		Config:          cfg,
	}
}
