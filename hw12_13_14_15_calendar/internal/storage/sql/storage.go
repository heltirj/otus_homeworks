package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/heltirj/otus_homeworks/hw12_13_14_15_calendar/internal/storage"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, dsn string) (*Storage, error) {
	strg := &Storage{}

	if err := strg.Connect(ctx, dsn); err != nil {
		return nil, err
	}

	return strg, nil
}

func (s *Storage) Connect(ctx context.Context, dsn string) error {
	var err error
	s.pool, err = pgxpool.Connect(ctx, dsn)
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	if err := s.pool.Ping(ctx); err != nil {
		return fmt.Errorf("error pinging database: %w", err)
	}

	return nil
}

func (s *Storage) Close() error {
	if s.pool == nil {
		return nil
	}

	s.pool.Close()
	s.pool = nil
	return nil
}

func (s *Storage) Add(ctx context.Context, event storage.Event) error {
	query := `INSERT INTO events (id, title, event_date, duration, description, user_id, notify_before) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := s.pool.Exec(ctx, query, event.ID, event.Title, event.EventDate, event.Duration,
		event.Description, event.UserID, event.NotifyBefore)
	return err
}

func (s *Storage) Update(ctx context.Context, updateEvent storage.UpdateEvent) error {
	query := `UPDATE events 
			  SET title = COALESCE($1, title), 
			      event_date = COALESCE($2, event_date), 
                  duration = COALESCE($3, duration), 
                  description = COALESCE($4, description), 
                  user_id = COALESCE($5, user_id), 
                  notify_before = COALESCE($6, notify_before) WHERE id = $7`
	_, err := s.pool.Exec(ctx, query,
		updateEvent.Title,
		updateEvent.EventDate,
		updateEvent.Duration,
		updateEvent.Description,
		updateEvent.UserID,
		updateEvent.NotifyBefore,
		updateEvent.ID)
	return err
}

func (s *Storage) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM events WHERE id = $1`
	_, err := s.pool.Exec(ctx, query, id)
	return err
}

func (s *Storage) GetByID(ctx context.Context, id uuid.UUID) (storage.Event, error) {
	query := `SELECT id, title, event_date, duration, description, user_id, notify_before FROM events WHERE id = $1`
	var event storage.Event
	err := s.pool.QueryRow(ctx, query, id).Scan(&event.ID, &event.Title, &event.EventDate,
		&event.Duration, &event.Description, &event.UserID, &event.NotifyBefore)
	if err != nil {
		return storage.Event{}, err
	}
	return event, nil
}

func (s *Storage) GetList(ctx context.Context, start time.Time, end time.Time) ([]storage.Event, error) {
	query := `SELECT id, title, event_date, duration, description, user_id, notify_before 
			  FROM events WHERE event_date >= $1 AND event_date < $2 ORDER BY event_date ASC`
	rows, err := s.pool.Query(ctx, query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []storage.Event
	for rows.Next() {
		var event storage.Event
		if err := rows.Scan(&event.ID, &event.Title, &event.EventDate, &event.Duration, &event.Description,
			&event.UserID, &event.NotifyBefore); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}
