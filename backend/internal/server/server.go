package server

import (
	"algoplatform/internal/config"
	httpi "algoplatform/internal/controller/http"
	"algoplatform/internal/controller/http/handlers"
	"algoplatform/internal/repo/postgres"
	"algoplatform/internal/usecase"
	"algoplatform/internal/worker"
	"algoplatform/pkg/db"
	"algoplatform/pkg/judge"
	"algoplatform/pkg/jwt"
	logger "algoplatform/pkg/log"
	"algoplatform/pkg/log/zap"
	"algoplatform/pkg/validator"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	serviceName  = "algoplatform"
	envFile      = ".env"
	AccessTTL    = 15 * time.Minute
	RefreshTTL   = 7 * 24 * time.Hour
)

func RunServer() {
	cfg, err := config.Load(envFile)
	if err != nil {
		log.Fatalf("Error load configs: %v", err)
	}

	logger, cleanup, err := zap.NewLogger(serviceName, cfg.Env)
	if err != nil {
		log.Fatalf("Error initialize logger: %v", err)
	}

	tokens := jwt.New(cfg.SecretKey, AccessTTL, RefreshTTL)

	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := db.NewDB(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Fatalf("Error initialize db: %v", err)
	}

	userRepo := postgres.NewUserRepo(db)
	problemRepo := postgres.NewProblemRepo(db)
	submissionRepo := postgres.NewSubmissionRepo(db)

	UserService := usecase.NewUserUsecase(userRepo)
	ProblemService := usecase.NewProblemUsecase(problemRepo)
	SubmissionService := usecase.NewSubmissionUsecase(submissionRepo)

	val := validator.New()

	UserHandler := handlers.NewUserHandler(UserService, tokens, val, logger)
	ProblemHandler := handlers.NewProblemHandler(ProblemService, val, logger)
	SubmissionHandler := handlers.NewSubmissionHandler(SubmissionService, val, logger)

	router := httpi.NewRouter(logger, UserHandler, ProblemHandler, SubmissionHandler, tokens)

	judgeClient := judge.NewClient(cfg.Judge0RapidAPIHost, cfg.Judge0RapidAPIKey)
	worker := worker.NewJudgeWorker(SubmissionService, ProblemService, judgeClient, logger)
	go worker.Start(ctx)

	go gracefulShutdown(cancel, logger)

	logger.Infof("Starting HTTP server on :%s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, router); err != nil {
		logger.Fatalf("Error running server: %v", err)
	}
}

func gracefulShutdown(cancel context.CancelFunc, logger logger.Logger) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	sig := <-sigChan
	logger.Infof("Received signal: %s. Shutting down gracefully...", sig)
	cancel()
}
