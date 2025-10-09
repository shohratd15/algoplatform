package httpi

import (
	"algoplatform/internal/controller/http/handlers"
	"algoplatform/internal/domain"
	"algoplatform/pkg/log"
	"fmt"
	"net/http"
)

func NewRouter(
	logger log.Logger,
	userHandler *handlers.UserHandler,
	problemHandler *handlers.ProblemHandler,
	submissionHandler *handlers.SubmissionHandler,
	tokenService domain.TokenService,
) *http.ServeMux {
	router := http.NewServeMux()

	auth := &AuthMiddleware{Tokens: tokenService}

	router.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprint(w, "Service is running! DB connection successful."); err != nil {
			logger.Errorf("ping error: %v", err)
		}
	})

	// Auth
	router.HandleFunc("POST /register", userHandler.Register)
	router.HandleFunc("POST /login", userHandler.Login)

	protected := http.NewServeMux()
	// Problems
	protected.HandleFunc("POST /problems", problemHandler.Create)
	protected.HandleFunc("GET /problems", problemHandler.List)
	protected.HandleFunc("GET /problems/detail", problemHandler.GetById)
	protected.HandleFunc("DELETE /problems", problemHandler.Delete)

	// Submissions
	protected.HandleFunc("POST /submissions", submissionHandler.Create)
	protected.HandleFunc("GET /submissions", submissionHandler.Get)

	protectedWithMiddleware := Logging(
		auth.JWT(RequireUser(protected)),
		logger,
	)

	mainRouter := http.NewServeMux()
	mainRouter.Handle("/", Logging(router, logger))
	mainRouter.Handle("/api/", http.StripPrefix("/api", protectedWithMiddleware))

	return mainRouter
}
