package storage

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID           uuid.UUID      `db:"id"`
	Title        string         `db:"title"`
	EventDate    time.Time      `db:"event_date"`
	Duration     time.Duration  `db:"duration"`
	Description  *string        `db:"description,omitempty"`
	UserID       int            `db:"user_id"`
	NotifyBefore *time.Duration `db:"notify_before,omitempty"`
}

type UpdateEvent struct {
	ID           uuid.UUID                         `db:"id"`
	Title        *UpdateNotNullable[string]        `db:"title,omitempty"`
	EventDate    *UpdateNotNullable[time.Time]     `db:"event_date,omitempty"`
	Duration     *UpdateNotNullable[time.Duration] `db:"duration,omitempty"`
	Description  *UpdateNullable[string]           `db:"description,omitempty"`
	UserID       *UpdateNotNullable[int]           `db:"user_id"`
	NotifyBefore *UpdateNullable[time.Duration]    `db:"notify_before,omitempty"`
}

type UpdateNotNullable[T any] struct {
	Value T
}

type UpdateNullable[T any] struct {
	Value *T
}

func NewUpdateNotNullable[T any](val T) *UpdateNotNullable[T] {
	return &UpdateNotNullable[T]{Value: val}
}

func NewUpdateNullable[T any](val *T) *UpdateNullable[T] {
	return &UpdateNullable[T]{Value: val}
}
