package di

import (
	"log"

	grpc "github.com/anfastk/mergespace/auth/internal/auth/adapter/inbound/grpc/handler"
	"github.com/anfastk/mergespace/auth/internal/auth/adapter/outbound/idgen"
	"github.com/anfastk/mergespace/auth/internal/auth/adapter/outbound/kafka"
	"github.com/anfastk/mergespace/auth/internal/auth/adapter/outbound/oauth"
	"github.com/anfastk/mergespace/auth/internal/auth/adapter/outbound/otp"
	"github.com/anfastk/mergespace/auth/internal/auth/adapter/outbound/postgres"
	"github.com/anfastk/mergespace/auth/internal/auth/adapter/outbound/redis"
	"github.com/anfastk/mergespace/auth/internal/auth/adapter/outbound/token"
	"github.com/anfastk/mergespace/auth/internal/auth/adapter/outbound/worker"
	"github.com/anfastk/mergespace/auth/internal/auth/application/port/inbound"
	"github.com/anfastk/mergespace/auth/internal/auth/application/usecase"
	"github.com/anfastk/mergespace/auth/internal/auth/infrastructure/config"
	"github.com/anfastk/mergespace/auth/internal/auth/infrastructure/crypto"
	"github.com/anfastk/mergespace/auth/internal/auth/infrastructure/database"
	redisconfig "github.com/anfastk/mergespace/auth/internal/auth/infrastructure/redisConfig"
	platformAvro "github.com/anfastk/mergespace/platform/infrastructure/avro"
	platformKafka "github.com/anfastk/mergespace/platform/infrastructure/kafka"
	"github.com/anfastk/mergespace/platform/infrastructure/messaging/schemas"
)

type App struct {
	Handler *grpc.AuthHandler

	HandlerUsecase inbound.AuthUseCase
	GoogleProvider *oauth.GoogleOAuthProvider
	GitHubProvider *oauth.GitHubOAuthProvider

	Worker *worker.OutboxWorker
}

func BuildApp() *App {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewPostgres(database.PostgresConfig{DSN: cfg.DatabaseURL})
	if err != nil {
		log.Fatalf("failed to init postgres: %v", err)
	}

	authRepo := postgres.NewUserRepository(db)
	idGen := idgen.NewUUIDGenerator()
	otpGen := otp.NewCryptoOTPGenerator()

	registry := platformAvro.NewRegistry(cfg.Kafka.SchemaRegistryURL)
	codec := platformAvro.NewCodec(registry)

	schemaBytes, err := schemas.FS.ReadFile("send_otp.avsc")
	if err != nil {
		log.Fatalf("failed to read send_otp.avsc: %v", err)
	}

	if err := codec.Register("auth.send_otp", "auth.notification-send_otp-value", string(schemaBytes)); err != nil {
		log.Fatalf("failed to register schema: %v", err)
	}

	kafkaProducer, err := platformKafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.UserSignupTopic, codec)
	if err != nil {
		log.Fatalf("kafka init failed: %v", err)
	}

	producer := kafka.NewEventProducer(kafkaProducer)

	redisClient := redisconfig.NewRedis(redisconfig.RedisConfig{
		Addr: cfg.Redis.Addr,
		DB:   0,
	})

	pendingSignupRepo := redis.NewSignupContextRedisStore(redisClient)
	passwordHash := crypto.NewBcryptHasher(16)

	outboxRepo := postgres.NewOutboxRepo(db)
	tokenGen := token.NewJWTGenerator(
		cfg.JWT.AccessSecret,
		cfg.JWT.RefreshSecret,
	)

	googleProvider := oauth.NewGoogleOAuthProvider(
		cfg.Google.ClientID,
		cfg.Google.ClientSecret,
		cfg.Google.RedirectURL,
	)

	githubProvider := oauth.NewGitHubOAuthProvider(
		cfg.GitHub.ClientID,
		cfg.GitHub.ClientSecret,
		cfg.GitHub.RedirectURL,
	)

	authIdentityRepo := postgres.NewAuthIdentityRepository(db)

	passwordResetStore := redis.NewPasswordResetRedisStore(redisClient)

	refreshSessionRepo := postgres.NewRefreshSessionRepo(db)

	authService := usecase.NewAuthService(
		db,
		authRepo,
		otpGen,
		idGen,
		pendingSignupRepo,
		passwordHash,
		producer,
		outboxRepo,
		tokenGen,
		passwordResetStore,
		authIdentityRepo,
		googleProvider,
		githubProvider,
		refreshSessionRepo,
	)

	handler := grpc.NewAuthHandler(authService)

	outboxWorker := worker.NewOutboxWorker(outboxRepo, producer)

	return &App{
		Handler:        handler,
		HandlerUsecase: authService,
		GoogleProvider: googleProvider,
		GitHubProvider: githubProvider,
		Worker:         outboxWorker,
	}
}
