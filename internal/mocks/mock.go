package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/amirhossein-karimi/todo/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockTodoRepository struct {
	mock.Mock
}

type MockCache struct {
	mock.Mock
}

func (m *MockCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	args := m.Called(ctx, key, value, ttl)
	return args.Error(0)
}

func (m *MockCache) Get(ctx context.Context, key string, dest interface{}) error {
	args := m.Called(ctx, key, dest)

	if args.Get(1) != nil {
		switch value := args.Get(1).(type) {
		case []*models.Todo:
			*(dest.(*[]*models.Todo)) = value

		case int64:
			*(dest.(*int64)) = value
		}
	}

	return args.Error(0)
}

func (m *MockCache) Delete(ctx context.Context, pattern string) error {
	args := m.Called(ctx, pattern)
	return args.Error(0)
}

func (m *MockTodoRepository) FindByTitle(ctx context.Context, title string) (*models.Todo, error) {
	args := m.Called(ctx, title)

	var todo *models.Todo

	if args.Get(0) != nil {
		todo = args.Get(0).(*models.Todo)
	}

	return todo, args.Error(1)
}

func (m *MockTodoRepository) FindByUUID(ctx context.Context, uuid uuid.UUID) (*models.Todo, error) {
	args := m.Called(ctx, uuid)

	var todo *models.Todo

	if args.Get(0) != nil {
		todo = args.Get(0).(*models.Todo)
	}

	return todo, args.Error(1)
}

func (m *MockTodoRepository) Delete(ctx context.Context, uuid uuid.UUID) error {
	args := m.Called(ctx, uuid)
	return args.Error(0)
}

func (m *MockTodoRepository) Create(ctx context.Context, todo *models.Todo) (uuid.UUID, error) {
	args := m.Called(ctx, todo)

	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockTodoRepository) List(ctx context.Context, page int, status string, priority string, assignee string) ([]*models.Todo, int64, error) {
	args := m.Called(ctx, page, status, priority, assignee)

	var todos []*models.Todo
	if args.Get(0) != nil {
		todos = args.Get(0).([]*models.Todo)
	}

	return todos, args.Get(1).(int64), args.Error(2)
}

func (m *MockTodoRepository) UpdateByUUID(ctx context.Context, uuid uuid.UUID, todo *models.Todo) (*models.Todo, error) {

	args := m.Called(ctx, uuid, todo)

	var updatedTodo *models.Todo

	if args.Get(0) != nil {
		updatedTodo = args.Get(0).(*models.Todo)
	}

	return updatedTodo, args.Error(1)
}

func (m *MockTodoRepository) FindByTitleExcludingUUID(ctx context.Context, uuid uuid.UUID, title string) (*models.Todo, error) {
	args := m.Called(ctx, uuid, title)

	var todo *models.Todo

	if args.Get(0) != nil {
		todo = args.Get(0).(*models.Todo)
	}

	return todo, args.Error(1)
}
