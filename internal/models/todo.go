package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	TodoStatusTodo = iota
	TodoStatusDoing
	TodoStatusDone
)

type Todo struct {
	ID   uint      `gorm:"primaryKey"`
	UUID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`

	Title       string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text"`

	Priority int `gorm:"not null;default:0"`
	Status   int `gorm:"type:integer;not null;default:0"`

	ToDoStartAt *time.Time `gorm:"type:timestamptz"`
	DoStartAt   *time.Time `gorm:"type:timestamptz"`
	DoneStartAt *time.Time `gorm:"type:timestamptz"`

	CreatedAt time.Time `gorm:"not null;default:now()"`
	UpdatedAt time.Time `gorm:"not null;default:now()"`
}

func (Todo) TableName() string {
	return "todos"
}
