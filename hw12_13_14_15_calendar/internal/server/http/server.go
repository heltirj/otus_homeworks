package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/heltirj/otus_homeworks/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	server *http.Server
	logger Logger
	app    Application
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	DebugKV(msg string, keysAndValues ...interface{})
	InfoKV(msg string, keysAndValues ...interface{})
	WarnKV(msg string, keysAndValues ...interface{})
	ErrorKV(msg string, keysAndValues ...interface{})
}

type Application interface {
	CreateEvent(ctx context.Context, title string, eventDate time.Time, duration time.Duration,
		description *string, userID int, notifyBefore *time.Duration) (uuid.UUID, error)
	UpdateEvent(ctx context.Context, updateEvent storage.UpdateEvent) error
	DeleteEvent(ctx context.Context, updateEvent storage.UpdateEvent) error
	GetDayEvents(ctx context.Context, dayDate time.Time) ([]storage.Event, error)
	GetWeekEvents(ctx context.Context, weekStart time.Time) ([]storage.Event, error)
	GetMonthEvents(ctx context.Context, weekStart time.Time) ([]storage.Event, error)
}

func NewServer(logger Logger, app Application, host string, port int) *Server {
	srv := &Server{
		logger: logger,
		app:    app,
	}

	srv.server = &http.Server{
		Addr:              fmt.Sprintf("%s:%d", host, port),
		Handler:           srv.setupRoutes(),
		ReadHeaderTimeout: time.Second * 10,
	}

	return srv
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error(fmt.Sprintf("Server failed: %s", err))
		}
	}()

	<-ctx.Done()

	return s.Stop(ctx)
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	return nil
}

func (s *Server) setupRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", s.loggingMiddleware(http.HandlerFunc(s.defaultHandler)))
	return mux
}

func (s *Server) defaultHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Events handler"))
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
