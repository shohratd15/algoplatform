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
) http.Handler {
	auth := &AuthMiddleware{Tokens: tokenService}

	// ── Ping (публичный, без /api префикса) ──────────────────────────────────
	publicRouter := http.NewServeMux()
	publicRouter.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprint(w, "Service is running! DB connection successful."); err != nil {
			logger.Errorf("ping error: %v", err)
		}
	})

	// ── Защищённые роуты для обычных пользователей ───────────────────────────
	userRouter := http.NewServeMux()
	userRouter.HandleFunc("GET /problems", problemHandler.List)
	userRouter.HandleFunc("GET /problems/detail", problemHandler.GetById)
	userRouter.HandleFunc("POST /submissions", submissionHandler.Create)
	userRouter.HandleFunc("GET /submissions", submissionHandler.Get)

	// ── Защищённые роуты для администраторов ─────────────────────────────────
	adminRouter := http.NewServeMux()
	adminRouter.HandleFunc("POST /problems", problemHandler.Create)
	adminRouter.HandleFunc("PUT /problems", problemHandler.Update)
	adminRouter.HandleFunc("DELETE /problems", problemHandler.Delete)

	// ── Главный роутер ────────────────────────────────────────────────────────
	mainRouter := http.NewServeMux()

	// Ping — без авторизации
	mainRouter.Handle("/", Logging(publicRouter, logger))

	// для auth — /api/register и /api/login.
	mainRouter.HandleFunc("POST /api/register", userHandler.Register)
	mainRouter.HandleFunc("POST /api/login", userHandler.Login)
	mainRouter.HandleFunc("POST /api/refresh", userHandler.Refresh)

	// Защищённые пользовательские роуты: /api/*
	mainRouter.Handle("/api/", http.StripPrefix("/api",
		Logging(auth.JWT(RequireUser(userRouter)), logger),
	))

	// Защищённые admin роуты: /api/admin/*
	mainRouter.Handle("/api/admin/", http.StripPrefix("/api/admin",
		Logging(auth.JWT(RequireAdmin(adminRouter)), logger),
	))

	return CORS(mainRouter)
}
