package httpi

import (
	"algoplatform/internal/controller/http/handlers"
	"algoplatform/pkg/log"
	"fmt"
	"net/http"
)

func NewRouter(
	logger log.Logger,
	userHandler *handlers.UserHandler,
	problemHandler *handlers.ProblemHandler,
) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprint(w, "Service is running! DB connection successful."); err != nil {
			logger.Errorf("ping error: %v", err)
		}
	})

	// Auth
	router.HandleFunc("POST /register", userHandler.Register)
	router.HandleFunc("POST /login", userHandler.Login)

	// Problems
	router.HandleFunc("POST /problems", problemHandler.Create)
	router.HandleFunc("GET /problems", problemHandler.List)
	router.HandleFunc("GET /problems/detail", problemHandler.GetBySlug)
	router.HandleFunc("DELETE /problems", problemHandler.Delete)

	return router
}
