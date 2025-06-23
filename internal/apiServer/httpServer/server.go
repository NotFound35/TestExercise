package httpServer

import (
	"awesomeProject/internal/apiServer/controllers"
	"awesomeProject/internal/config"
	"awesomeProject/internal/userservice"
	"context"
	"errors"
	"fmt"
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
}

func NewServer(u userservice.IUserService, log *zap.Logger) *Server {
	server := &Server{
		router:   chi.NewRouter(),
		handlers: controllers.NewHandler(u, log),
	}
	server.setupRoutes()
	return server
}

func (s *Server) Run(cfg *config.Config) {
	const op = "Run"
	fmt.Println("Cервер запущен")
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
		fmt.Println("Ожидание запросов")
		err := http.ListenAndServe(cfg.HTTPServer.Address, s.router)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Errorf("метод %v: %v", op, err)
		}
	}()

	<-stop
	fmt.Println("Сервер выключается...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		fmt.Errorf("метод %v: %v", op, err)
	}
	fmt.Println("Сервер завершен КОРРЕКТНО")
}

func (s *Server) Middleware(next http.Handler) http.Handler {
	const op = "Middleware"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				//todo log / return
				fmt.Errorf("метод %v: %v", op, err)
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
		_, err := w.Write([]byte("Сервер работает")) //проверка работоспособности сервера
		if err != nil {
			fmt.Errorf("метод %v: %v", op, err)
		}
	})

	s.router.Post("/users", s.handlers.SaveUserHandler)
	s.router.Get("/users/search", s.handlers.GetUserHandler)
	s.router.Get("/users/list", s.handlers.ListUsersHandler)
	s.router.Delete("/users/{id}", s.handlers.DeleteUserHandler)
	s.router.Delete("/users/soft/{id}", s.handlers.SoftDeleteUserHandler)
	s.router.Patch("/users/update/{id}", s.handlers.UpdateUserHandler)

}
