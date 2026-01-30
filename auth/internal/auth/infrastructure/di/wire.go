package di

import (
	"log"

	grpc "github.com/anfastk/mergespace/auth/internal/auth/adapter/inbound/grpc/handler"
	"github.com/anfastk/mergespace/auth/internal/auth/adapter/outbound/idgen"
	"github.com/anfastk/mergespace/auth/internal/auth/adapter/outbound/kafka"
	"github.com/anfastk/mergespace/auth/internal/auth/adapter/outbound/otp"
	"github.com/anfastk/mergespace/auth/internal/auth/adapter/outbound/postgres"
	"github.com/anfastk/mergespace/auth/internal/auth/adapter/outbound/redis"
	"github.com/anfastk/mergespace/auth/internal/auth/application/service"
	"github.com/anfastk/mergespace/auth/internal/auth/infrastructure/config"
	"github.com/anfastk/mergespace/auth/internal/auth/infrastructure/crypto"
	"github.com/anfastk/mergespace/auth/internal/auth/infrastructure/database"
	redisconfig "github.com/anfastk/mergespace/auth/internal/auth/infrastructure/redisConfig"
	platformAvro "github.com/anfastk/mergespace/platform/infrastructure/avro"
	platformKafka "github.com/anfastk/mergespace/platform/infrastructure/kafka"
	"github.com/anfastk/mergespace/platform/infrastructure/messaging/schemas"
)

func BuildAuthHandler() *grpc.AuthHandler {
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
	usernameGen := service.NewUsernameAllocator(authRepo, 10)

	registry := platformAvro.NewRegistry(cfg.Kafka.SchemaRegistryURL)

	codec := platformAvro.NewCodec(registry)

	schemaBytes, err := schemas.FS.ReadFile("send_otp.avsc")
	if err != nil {
		log.Fatalf("failed to read send_otp.avsc: %v", err)
	}

	if err := codec.Register("auth.send_otp", "auth.notification-send_otp-value", string(schemaBytes)); err != nil {
		log.Fatalf("failed to register auth.send_otp schema: %v", err)
	}

	kafkaProducer, err := platformKafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.UserSignupTopic, codec)

	if err != nil {
		log.Fatalf("schema registration failed: %v", err)
	}

	producer := kafka.NewEventProducer(kafkaProducer)

	redisClient := redisconfig.NewRedis(redisconfig.RedisConfig{
		Addr: cfg.Redis.Addr,
		DB:   0,
	})

	pendingSignupRepo := redis.NewSignupContextRedisStore(redisClient)
	passwordHash := crypto.NewBcryptHasher(16)

	authService := service.NewAuthService(
		authRepo,
		usernameGen,
		otpGen,
		idGen,
		pendingSignupRepo,
		passwordHash,
		producer,
	)

	return grpc.NewAuthHandler(authService)
}
