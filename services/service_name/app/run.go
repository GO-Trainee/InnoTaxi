package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(configPath string) error {
	// -------------------------------------------------------------------------
	// Observability: structured logger (slog).
	// Replace slog.NewJSONHandler with slog.NewTextHandler for local dev.
	// In production wire in your OTel exporter here instead.
	// -------------------------------------------------------------------------
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// -------------------------------------------------------------------------
	// Root context — cancelled on OS signal (SIGTERM / SIGINT).
	// Everything downstream receives this context and must respect its Done().
	// -------------------------------------------------------------------------
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// -------------------------------------------------------------------------
	// Configuration
	// -------------------------------------------------------------------------
	// cfg, err := config.Load(configPath)
	// if err != nil {
	//     return fmt.Errorf("load config: %w", err)
	// }

	// -------------------------------------------------------------------------
	// Database connections + migrations
	// -------------------------------------------------------------------------
	// Migrations run inside New() — if they fail, the call returns an error
	// and the service must not start. Do NOT call migrations separately here.
	// pgClient, err  := pg.New(cfg.Postgres)
	// if err != nil {
	//     return fmt.Errorf("pg init: %w", err)
	// }
	// mongoClient, err := mongo.New(cfg.Mongo)
	// if err != nil {
	//     return fmt.Errorf("mongo init: %w", err)
	// }
	// redisClient := redis.New(cfg.Redis)

	// -------------------------------------------------------------------------
	// Infrastructure: gRPC client, Kafka producer
	// -------------------------------------------------------------------------
	// grpcClient, err := appgrpc.NewClient(cfg.GRPC)
	// if err != nil {
	//     return fmt.Errorf("grpc client: %w", err)
	// }
	// kafkaPublisher := appkafka.NewPublisher(cfg.Kafka)

	// -------------------------------------------------------------------------
	// Repositories
	// -------------------------------------------------------------------------
	// pgRepo    := pg.NewPgRepository(pgClient.DB())
	// mongoRepo := mongo.NewMongoRepository(mongoClient)
	// redisRepo := redis.New(redisClient)

	// -------------------------------------------------------------------------
	// Gateways (outgoing calls)
	// -------------------------------------------------------------------------
	// grpcGateway  := gatewaygrpc.New(grpcClient)
	// httpGateway  := gatewayhttp.New(cfg.ThirdPartyURL)
	// kafkaGateway := gatewaykafka.New(kafkaPublisher)

	// -------------------------------------------------------------------------
	// Services (business logic)
	// -------------------------------------------------------------------------
	// svc := service.New(pgRepo, grpcGateway, kafkaGateway)

	// -------------------------------------------------------------------------
	// Handlers (incoming requests)
	// -------------------------------------------------------------------------
	// httpHandler  := handlerhttp.New(svc)
	// grpcHandler  := handlergrpc.New(svc)
	// kafkaHandler := handlerkafka.New(svc)

	// -------------------------------------------------------------------------
	// Servers — start in goroutines, collect errors via errCh.
	// -------------------------------------------------------------------------
	errCh := make(chan error, 3)

	// HTTP server
	httpServer := &http.Server{
		Addr: ":8080", // cfg.HTTP.Addr
		// Handler: httpHandler.Router(),
	}
	go func() {
		logger.Info("http server starting", "addr", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	// gRPC server
	// grpcSrv, err := appgrpc.NewServer(grpcHandler)
	// if err != nil {
	//     return fmt.Errorf("grpc server init: %w", err)
	// }
	// go func() {
	//     logger.Info("grpc server starting", "addr", cfg.GRPC.Addr)
	//     if err := grpcSrv.Serve(); err != nil {
	//         errCh <- err
	//     }
	// }()

	// Kafka consumer
	// kafkaConsumer := appkafka.NewConsumer(cfg.Kafka, kafkaHandler)
	// go func() {
	//     logger.Info("kafka consumer starting")
	//     if err := kafkaConsumer.Run(ctx); err != nil {
	//         errCh <- err
	//     }
	// }()

	// -------------------------------------------------------------------------
	// Block until shutdown signal or a server error.
	// -------------------------------------------------------------------------
	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	case err := <-errCh:
		logger.Error("server error", "err", err)
	}

	// -------------------------------------------------------------------------
	// Graceful shutdown — give servers 15 s to drain in-flight requests.
	// -------------------------------------------------------------------------
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	logger.Info("shutting down http server")
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("http server shutdown error", "err", err)
	}

	// grpcSrv.GracefulStop()

	// Flush kafka producer before exit so no messages are lost.
	// kafkaPublisher.Close()

	logger.Info("shutdown complete")
	return nil
}
