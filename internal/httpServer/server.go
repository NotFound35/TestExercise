package httpServer

import (
	"awesomeProject/config"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// структура сервера
type Server struct {
	httpServer *http.Server
	config     *config.Config
}

// создание нового сервера с настройками из конфига
func New(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
		httpServer: &http.Server{
			Addr:         cfg.Server.Address,     //адрес
			ReadTimeout:  cfg.Server.Timeout,     //таймаут чтения
			WriteTimeout: cfg.Server.Timeout,     //таймаут записи
			IdleTimeout:  cfg.Server.IdleTimeout, //таймаут бездействия
		},
	}
}

// запуск httpServer
func (s *Server) Start(handler http.Handler) {
	s.httpServer.Handler = s.recoveryMiddleware(handler) //добавили обработчика паник

	go func() { //запуск самого сервера
		log.Println("HTTP сервер запущен")
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln("Ошибка сервера")
		}

		s.gracefulShutdown()
	}()
}

// плавное завершение работы - говорили на уроке
func (s *Server) gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt) //ловим сигнал ctrl+c

	<-quit //ждем сигнал
	log.Println("Завершение работы сервера")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Stop(ctx); err != nil { //пытаемся остановиться
		log.Println("Ошибка при завершении работы сервера")
	}

	log.Println("Сервер успешно остановлен")
}

// остановка сервера - используем выше для плавной остановки сервера
func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// перехватывает паники в обработчиках - для того, что бы сервер не падал при ошибках в обработчиках
func (s *Server) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil { //ловим панику
				log.Println("Перехвачена паника")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("Внутренняя ошибка сервера"))
			}
		}()
		next.ServeHTTP(w, r) //продолжаем цепочку обработчиков
	})
}
