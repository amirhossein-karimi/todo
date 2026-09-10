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
	ID   uint      `gorm:"primaryKey" json:"-"`
	UUID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"uuid"`

	Assignee    string `gorm:"type:varchar(255);not null" json:"assignee"`
	Title       string `gorm:"type:varchar(255);not null" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	Priority    int    `gorm:"not null;default:1" json:"priority"`
	Status      int    `gorm:"type:integer;not null;default:0" json:"status"`

	ToDoStartAt *time.Time `gorm:"type:timestamptz" json:"todo_start_at"`
	DoStartAt   *time.Time `gorm:"type:timestamptz" json:"do_start_at"`
	DoneStartAt *time.Time `gorm:"type:timestamptz" json:"done_start_at"`

	CreatedAt time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

func (Todo) TableName() string {
	return "todos"
}
