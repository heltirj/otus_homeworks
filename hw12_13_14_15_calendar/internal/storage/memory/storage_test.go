package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/heltirj/otus_homeworks/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestAddEvent(t *testing.T) {
	strg := New()
	eventID := uuid.New()
	event := storage.Event{
		ID:           eventID,
		Title:        "Test Event",
		EventDate:    time.Now().Add(1 * time.Hour),
		Duration:     time.Hour,
		Description:  nil,
		UserID:       1,
		NotifyBefore: nil,
	}

	ctx := context.Background()
	if err := strg.Add(ctx, event); err != nil {
		t.Errorf("Expected no error adding event, got %v", err)
	}

	retrievedEvent, err := strg.GetByID(ctx, eventID)
	if err != nil {
		t.Errorf("Expected no error getting event, got %v", err)
	}

	if retrievedEvent.ID != eventID {
		t.Errorf("Expected event ID %v, got %v", eventID, retrievedEvent.ID)
	}
}

func TestUpdateEvent(t *testing.T) {
	strg := New()
	eventID := uuid.New()
	event := storage.Event{
		ID:           eventID,
		Title:        "Initial Event",
		EventDate:    time.Now().Add(1 * time.Hour),
		Duration:     time.Hour,
		Description:  nil,
		UserID:       1,
		NotifyBefore: nil,
	}

	ctx := context.Background()

	if err := strg.Add(ctx, event); err != nil {
		t.Fatalf("Failed to add event: %v", err)
	}

	newTitle := "Updated Event"
	updateEvent := storage.UpdateEvent{
		ID:           eventID,
		Title:        storage.NewUpdateNotNullable(newTitle),
		EventDate:    nil,
		Duration:     nil,
		Description:  nil,
		NotifyBefore: nil,
	}

	if err := strg.Update(ctx, updateEvent); err != nil {
		t.Errorf("Expected no error updating event, got %v", err)
	}

	retrievedEvent, err := strg.GetByID(ctx, eventID)
	if err != nil {
		t.Errorf("Expected no error getting updated event, got %v", err)
	}

	if retrievedEvent.Title != newTitle {
		t.Errorf("Expected event title %q, got %q", newTitle, retrievedEvent.Title)
	}
}

func TestDeleteEvent(t *testing.T) {
	strg := New()
	eventID := uuid.New()
	event := storage.Event{
		ID:           eventID,
		Title:        "Event to Delete",
		EventDate:    time.Now().Add(1 * time.Hour),
		Duration:     time.Hour,
		Description:  nil,
		UserID:       1,
		NotifyBefore: nil,
	}

	ctx := context.Background()
	if err := strg.Add(ctx, event); err != nil {
		t.Fatalf("Failed to add event: %v", err)
	}

	if err := strg.Delete(ctx, eventID); err != nil {
		t.Errorf("Expected no error deleting event, got %v", err)
	}

	_, err := strg.GetByID(ctx, eventID)
	if err == nil {
		t.Errorf("Expected error getting deleted event, got none")
	}
}

func TestGetDayEvents(t *testing.T) {
	t.Parallel()
	strg := New()

	now := time.Now()
	event1 := storage.Event{
		ID:           uuid.New(),
		Title:        "Day Event 1",
		EventDate:    now.Truncate(24 * time.Hour).Add(1 * time.Hour),
		Duration:     time.Hour,
		Description:  nil,
		UserID:       1,
		NotifyBefore: nil,
	}

	event2 := storage.Event{
		ID:           uuid.New(),
		Title:        "Day Event 2",
		EventDate:    now.Truncate(24 * time.Hour).Add(2 * time.Hour),
		Duration:     time.Hour,
		Description:  nil,
		UserID:       1,
		NotifyBefore: nil,
	}

	event3 := storage.Event{
		ID:           uuid.New(),
		Title:        "Another Day Event",
		EventDate:    now.Truncate(24 * time.Hour).Add(25 * time.Hour),
		Duration:     time.Hour,
		Description:  nil,
		UserID:       1,
		NotifyBefore: nil,
	}

	ctx := context.Background()
	if err := strg.Add(ctx, event1); err != nil {
		t.Fatalf("Failed to add event1: %v", err)
	}
	if err := strg.Add(ctx, event2); err != nil {
		t.Fatalf("Failed to add event2: %v", err)
	}
	if err := strg.Add(ctx, event3); err != nil {
		t.Fatalf("Failed to add event3: %v", err)
	}

	today := time.Now().Truncate(24 * time.Hour)
	events, err := strg.GetList(ctx, today.Truncate(time.Hour*24), today.AddDate(0, 0, 1))
	if err != nil {
		t.Errorf("Expected no error getting day's events, got %v", err)
	}

	require.Equal(t, []storage.Event{event1, event2}, events)

	if len(events) != 2 {
		t.Errorf("Expected 2 events, got %d", len(events))
	}
}

func TestGetWeekEvents(t *testing.T) {
	t.Parallel()

	strg := New()

	weekStart := time.Date(2024, 11, 25, 0, 0, 0, 0, time.UTC)

	event1 := storage.Event{
		ID:           uuid.New(),
		Title:        "Week Event 1",
		EventDate:    weekStart.Add(1 * time.Hour),
		Duration:     time.Hour,
		Description:  nil,
		UserID:       1,
		NotifyBefore: nil,
	}

	event2 := storage.Event{
		ID:           uuid.New(),
		Title:        "Week Event 2",
		EventDate:    weekStart.Add(3 * 24 * time.Hour),
		Duration:     time.Hour,
		Description:  nil,
		UserID:       1,
		NotifyBefore: nil,
	}

	ctx := context.Background()
	if err := strg.Add(ctx, event1); err != nil {
		t.Fatalf("Failed to add event1: %v", err)
	}
	if err := strg.Add(ctx, event2); err != nil {
		t.Fatalf("Failed to add event2: %v", err)
	}

	events, err := strg.GetList(ctx, weekStart, weekStart.AddDate(0, 0, 7))
	if err != nil {
		t.Errorf("Expected no error getting week's events, got %v", err)
	}

	if len(events) != 2 {
		t.Errorf("Expected 2 events this week, got %d", len(events))
	}

	require.Equal(t, []storage.Event{event1, event2}, events)
}

func TestGetMonthEvents(t *testing.T) {
	t.Parallel()

	strg := New()

	monthStart := time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)

	event1 := storage.Event{
		ID:           uuid.New(),
		Title:        "Month Event 1",
		EventDate:    monthStart.Add(1 * 24 * time.Hour),
		Duration:     time.Hour,
		Description:  nil,
		UserID:       1,
		NotifyBefore: nil,
	}

	event2 := storage.Event{
		ID:           uuid.New(),
		Title:        "Month Event 2",
		EventDate:    monthStart.Add(30 * 24 * time.Hour),
		Duration:     time.Hour,
		Description:  nil,
		UserID:       1,
		NotifyBefore: nil,
	}

	ctx := context.Background()
	if err := strg.Add(ctx, event1); err != nil {
		t.Fatalf("Failed to add event1: %v", err)
	}
	if err := strg.Add(ctx, event2); err != nil {
		t.Fatalf("Failed to add event2: %v", err)
	}

	events, err := strg.GetList(ctx, monthStart, monthStart.AddDate(0, 1, 0))
	if err != nil {
		t.Errorf("Expected no error getting month's events, got %v", err)
	}

	if len(events) != 1 {
		t.Errorf("Expected 1 event this month, got %d", len(events))
	}

	require.Equal(t, []storage.Event{event1}, events)
}
