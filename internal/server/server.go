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
	"algoplatform/pkg/log/zap"
	"context"
	"log"
	"net/http"
	"time"
)

const (
	serviceName       = "algoplatform"
	envFile           = ".env"
	AppContextTimeout = 500 * time.Second
	SecretKeyTTL      = 48 * time.Hour
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

	tokens := jwt.New(cfg.SecretKey, SecretKeyTTL)

	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), AppContextTimeout)
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

	UserHandler := handlers.NewUserHandler(UserService, tokens, logger)
	ProblemHandler := handlers.NewProblemHandler(ProblemService, logger)
	SubmissionHandler := handlers.NewSubmissionHandler(SubmissionService, logger)

	router := httpi.NewRouter(logger, UserHandler, ProblemHandler, SubmissionHandler, tokens)

	judgeClient := judge.NewClient("http://judge0:2358")
	worker := worker.NewJudgeWorker(SubmissionService, ProblemService, judgeClient, logger)
	go worker.Start(ctx)

	logger.Info("Starting HTTP server on :8080")
	if err := http.ListenAndServe(":"+cfg.ServerPort, router); err != nil {
		logger.Fatalf("Error running server: %v", err)
	}
}
