package app

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/heltirj/otus_homeworks/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	Logger  Logger
	Storage Storage
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

type Storage interface {
	Add(ctx context.Context, event storage.Event) error
	Update(ctx context.Context, updateEvent storage.UpdateEvent) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (storage.Event, error)
	GetList(ctx context.Context, start time.Time, end time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, title string, eventDate time.Time, duration time.Duration,
	description *string, userID int, notifyBefore *time.Duration,
) (uuid.UUID, error) {
	id := uuid.New()
	return id, a.Storage.Add(ctx, storage.Event{
		ID: id, Title: title, EventDate: eventDate, Duration: duration,
		Description: description, UserID: userID, NotifyBefore: notifyBefore,
	})
}

func (a *App) UpdateEvent(ctx context.Context, updateEvent storage.UpdateEvent) error {
	return a.Storage.Update(ctx, updateEvent)
}

func (a *App) DeleteEvent(ctx context.Context, updateEvent storage.UpdateEvent) error {
	return a.Storage.Update(ctx, updateEvent)
}

func (a *App) GetDayEvents(ctx context.Context, dayDate time.Time) ([]storage.Event, error) {
	return a.Storage.GetList(ctx, dayDate.Truncate(time.Hour*24), dayDate.AddDate(0, 0, 1))
}

func (a *App) GetWeekEvents(ctx context.Context, weekStart time.Time) ([]storage.Event, error) {
	return a.Storage.GetList(ctx, weekStart, weekStart.AddDate(0, 0, 7))
}

func (a *App) GetMonthEvents(ctx context.Context, monthStart time.Time) ([]storage.Event, error) {
	return a.Storage.GetList(ctx, monthStart, monthStart.AddDate(0, 1, 0))
}
