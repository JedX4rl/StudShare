package main

import (
	mymongo "StudShare/internal/adapters/mongo"
	"StudShare/internal/adapters/postgres"
	myredis "StudShare/internal/adapters/redis"
	mys3 "StudShare/internal/adapters/s3"
	"StudShare/internal/auth"
	"StudShare/internal/config/server_config"
	"StudShare/internal/config/storage_config"
	"StudShare/internal/repository"
	"StudShare/internal/router"
	"StudShare/internal/router/handler"
	"StudShare/internal/service"
	"context"
	"errors"
	"github.com/joho/godotenv"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("Error loading .env file")
		os.Exit(1)
	}
}

// @title           StudShare API
// @version         1.0
// @description     API для виртуальной доски объявлений между студентами
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	serverCfg, err := server_config.MustLoadServerConfig()
	if err != nil {
		slog.Error("Error loading server config")
		os.Exit(1)
	}

	storageCfg, err := storage_config.MustLoadStorageConfig()
	if err != nil {
		slog.Error("Error loading storage config")
		os.Exit(1)
	}

	postgresDb, err := postgres.NewStorage(storageCfg.Postgres)
	if err != nil {
		slog.Error("failed to connect to postgres", err)
		os.Exit(1)
	}

	redis, err := myredis.NewRedis(storageCfg.Redis)
	if err != nil {
		slog.Error("failed to connect to redis", err)
		os.Exit(1)
	}

	mongo, err := mymongo.NewMongo(storageCfg.Mongo)
	if err != nil {
		slog.Error("failed to connect to mongodb", err)
		os.Exit(1)
	}

	s3, err := mys3.NewS3Client(storageCfg.S3)
	if err != nil {
		slog.Error("failed to connect to s3", err)
		os.Exit(1)
	}

	slog.Info("connected to database")

	defer func() {
		err = postgresDb.Close()
		if err != nil {
			slog.Error("got error when closing the DB connection", err)
			os.Exit(1)
		}
	}()

	repos := repository.NewRepository(postgresDb, redis, mongo, s3, storageCfg.S3)

	tokenManager, err := auth.NewTokenManager(repos.CacheRepo)
	if err != nil {
		slog.Error("failed to initialize token manager", err)
		os.Exit(1)
	}

	services := service.NewService(repos, tokenManager)

	handlers := handler.NewHandler(services, tokenManager)

	slog.Info("starting server")

	server := &http.Server{
		Addr:    serverCfg.Address,
		Handler: router.NewRouter(handlers),
	}

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("starting server...")
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("listen and serve error", err)
		}
	}()

	<-quit
	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", err)
	}

	slog.Info("server exiting")
}
