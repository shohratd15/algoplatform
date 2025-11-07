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

	// Auth (Public)
	router.HandleFunc("POST /register", userHandler.Register)
	router.HandleFunc("POST /login", userHandler.Login)

	protectedUserRouter := http.NewServeMux()

	// Problems
	protectedUserRouter.HandleFunc("GET /problems", problemHandler.List)
	protectedUserRouter.HandleFunc("GET /problems/detail", problemHandler.GetById)

	// Submissions
	protectedUserRouter.HandleFunc("POST /submissions", submissionHandler.Create)
	protectedUserRouter.HandleFunc("GET /submissions", submissionHandler.Get)

	protectedUserMiddleware := Logging(
		auth.JWT(RequireUser(protectedUserRouter)),
		logger,
	)

	adminRouter := http.NewServeMux()

	// Problems
	adminRouter.HandleFunc("POST /problems", problemHandler.Create)
	adminRouter.HandleFunc("DELETE /problems", problemHandler.Delete)

	protectedAdminMiddleware := Logging(
		auth.JWT(RequireAdmin(adminRouter)),
		logger,
	)

	mainRouter := http.NewServeMux()

	mainRouter.Handle("/", Logging(router, logger))

	mainRouter.Handle("/api/", http.StripPrefix("/api", protectedUserMiddleware))

	mainRouter.Handle("/api/admin/", http.StripPrefix("/api/admin", protectedAdminMiddleware))

	return mainRouter
}
