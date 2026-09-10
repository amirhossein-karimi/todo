package service

import (
	"context"
	"testing"
	"time"

	realCache "github.com/amirhossein-karimi/todo/internal/cache"
	"github.com/amirhossein-karimi/todo/internal/errors"
	"github.com/amirhossein-karimi/todo/internal/models"
	"github.com/amirhossein-karimi/todo/internal/requests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type mockTodoRepository struct {
	mock.Mock
}

type mockCache struct {
	mock.Mock
}

func (m *mockCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	args := m.Called(ctx, key, value, ttl)
	return args.Error(0)
}

func (m *mockCache) Get(ctx context.Context, key string, dest interface{}) error {
	args := m.Called(ctx, key, dest)
	return args.Error(0)
}

func (m *mockCache) Delete(ctx context.Context, pattern string) error {
	args := m.Called(ctx, pattern)
	return args.Error(0)
}

func (m *mockTodoRepository) FindByTitle(
	ctx context.Context,
	title string,
) (*models.Todo, error) {
	args := m.Called(ctx, title)

	var todo *models.Todo

	if args.Get(0) != nil {
		todo = args.Get(0).(*models.Todo)
	}

	return todo, args.Error(1)
}

func (m *mockTodoRepository) Create(
	ctx context.Context,
	todo *models.Todo,
) (uuid.UUID, error) {
	args := m.Called(ctx, todo)

	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *mockTodoRepository) List(ctx context.Context, page int) ([]*models.Todo, int64, error) {
	args := m.Called(ctx, page)

	var todos []*models.Todo
	if args.Get(0) != nil {
		todos = args.Get(0).([]*models.Todo)
	}

	return todos, args.Get(1).(int64), args.Error(2)
}

func TestTodoService_Create(t *testing.T) {
	ctx := context.Background()

	// arrange
	tests := []struct {
		name              string
		req               *requests.TodoRequest
		findByTitleResult *models.Todo
		findByTitleErr    error
		createResult      uuid.UUID
		createErr         error
		expectedUUID      uuid.UUID
		expectedErr       error
	}{
		{
			name: "should create todo successfully",
			req: &requests.TodoRequest{
				Title:       "Learn Go",
				Description: "Learn Go testing",
				Priority:    1,
			},
			findByTitleResult: nil,
			findByTitleErr:    gorm.ErrRecordNotFound,
			createResult:      uuid.New(),
			createErr:         nil,
		},
		{
			name: "should return error when title already exists",
			req: &requests.TodoRequest{
				Title:       "Learn Go",
				Description: "Testing Go",
				Priority:    1,
			},
			findByTitleResult: &models.Todo{
				UUID:  uuid.New(),
				Title: "Learn Go",
			},
			findByTitleErr: nil,
			expectedUUID:   uuid.Nil,
			expectedErr:    errors.ErrTodoWithThisTitleAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mockTodoRepository)
			cache := new(mockCache)
			service := NewTodoService(repo, cache)

			repo.
				On("FindByTitle", ctx, tt.req.Title).
				Return(tt.findByTitleResult, tt.findByTitleErr)

			realCache := realCache.CreateKeyName(realCache.TodosKey, "*")

			cache.On("Delete", ctx, realCache).Return(nil)

			if tt.findByTitleResult == nil {
				repo.
					On(
						"Create",
						ctx,
						mock.AnythingOfType("*models.Todo"),
					).
					Return(tt.createResult, tt.createErr)
			}

			// act
			result, err := service.Create(ctx, tt.req)

			// asserets
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
				assert.Equal(t, uuid.Nil, tt.expectedUUID)

				repo.AssertNotCalled(
					t,
					"Create",
					mock.Anything,
					mock.Anything,
				)
				cache.AssertNotCalled(t, "Delete", mock.Anything)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.createResult, result)

			repo.AssertExpectations(t)
		})
	}
}
