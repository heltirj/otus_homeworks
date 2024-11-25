package memorystorage

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/heltirj/otus_homeworks/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu     sync.RWMutex
	events map[uuid.UUID]storage.Event
}

func New() *Storage {
	return &Storage{
		events: make(map[uuid.UUID]storage.Event),
	}
}

func (s *Storage) Add(_ context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events[event.ID] = event
	return nil
}

func (s *Storage) Update(_ context.Context, updateEvent storage.UpdateEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	event, exists := s.events[updateEvent.ID]
	if !exists {
		return fmt.Errorf("event not found")
	}

	if updateEvent.Title != nil {
		event.Title = updateEvent.Title.Value
	}
	if updateEvent.EventDate != nil {
		event.EventDate = updateEvent.EventDate.Value
	}
	if updateEvent.Duration != nil {
		event.Duration = updateEvent.Duration.Value
	}
	if updateEvent.Description != nil {
		event.Description = updateEvent.Description.Value
	}
	if updateEvent.NotifyBefore != nil {
		event.NotifyBefore = updateEvent.NotifyBefore.Value
	}

	s.events[updateEvent.ID] = event
	return nil
}

func (s *Storage) Delete(_ context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.events[id]; !exists {
		return fmt.Errorf("event not found")
	}
	delete(s.events, id)
	return nil
}

func (s *Storage) GetByID(_ context.Context, id uuid.UUID) (storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	event, exists := s.events[id]
	if !exists {
		return storage.Event{}, fmt.Errorf("event not found")
	}
	return event, nil
}

func (s *Storage) GetList(_ context.Context, start time.Time, end time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var events []storage.Event
	for _, event := range s.events {
		if event.EventDate.After(start) && event.EventDate.Before(end) {
			events = append(events, event)
		}
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].EventDate.Before(events[j].EventDate)
	})

	return events, nil
}
