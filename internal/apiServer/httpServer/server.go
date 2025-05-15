package httpServer

import (
	"awesomeProject/internal/apiServer/controller"
	"awesomeProject/internal/config"
	"awesomeProject/internal/userservice"
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	router   *chi.Mux //TODO what is it?
	handlers *controller.Handler
}

// init server
func NewServer(u *userservice.UserService) *Server {
	server := &Server{
		router:   chi.NewRouter(),
		handlers: controller.NewHandler(u),
	}
	server.setupRoutes()
	return server
}

func (s *Server) StartAndFinish(cfg *config.Config) {
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
			log.Fatalln("Критическая ошибка сервера")
		}
	}()

	<-stop
	fmt.Println("Сервер выключается...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		log.Fatalln("Сервер завершен НЕ корректно")
	}
	fmt.Println("Сервер завершен КОРРЕКТНО")
}

func (s *Server) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("восстановлено после паники")
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *Server) setupRoutes() {
	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Сервер работает")) //проверка работоспособности сервера
		if err != nil {
			log.Fatalln("Миша, все фигня, давай по новой")
		}
	})
	s.router.Post("/users", s.handlers.SaveUserHandler) //регистрация обработчика для POST-запросов
}

/*Запись ответа:
w.Write отправляет клиенту сырые байты
Конвертирует строку "Сервер работат" в байтовый срез ([]byte)
Возвращает количество записанных байт и ошибку (количество игнорируется через _)*/
