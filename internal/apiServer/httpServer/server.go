package httpServer

import (
	"awesomeProject/internal/apiServer/controllers"
	"awesomeProject/internal/config"
	"awesomeProject/internal/userservice"
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Server struct {
	router   *chi.Mux
	handlers *controllers.Handler
	log      *zap.Logger
}

func NewServer(u userservice.IUserService, log *zap.Logger) *Server {
	server := &Server{
		router:   chi.NewRouter(),
		handlers: controllers.NewHandler(u, log),
		log:      log,
	}
	server.setupRoutes()
	return server
}

func (s *Server) Run(cfg *config.Config) {
	const op = "Server.Run"
	s.log.Info("server is starting")
	server := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      s.router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		s.log.Info("waiting request")
		err := http.ListenAndServe(cfg.HTTPServer.Address, s.router)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Fatal("error with starting server",
				zap.String("op", op),
				zap.Error(err),
			)
		}
	}()

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		s.log.Error(" ",
			zap.String("op", op),
			zap.Error(err),
		)
	}
}

func (s *Server) Middleware(next http.Handler) http.Handler {
	const op = "Middleware"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				s.log.Error("error with Middleware",
					zap.String("op", op),
					zap.Error(err.(error)),
				)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *Server) setupRoutes() {
	const op = "setupRoutes"
	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Server is running"))
		if err != nil {
			s.log.Error(" ",
				zap.String("op", op),
				zap.Error(err),
			)
		}
	})

	s.router.Post("/users", s.handlers.SaveUserHandler)
	s.router.Get("/users/search", s.handlers.GetUserHandler)
	s.router.Get("/users/list", s.handlers.ListUsersHandler)
	s.router.Delete("/users/{id}", s.handlers.DeleteUserHandler)
	s.router.Delete("/users/{id}/soft", s.handlers.SoftDeleteUserHandler)
	s.router.Patch("/users/{id}/update", s.handlers.UpdateUserHandler)

}
