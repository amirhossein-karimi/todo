package service

import (
	"context"
	"testing"
	"time"

	realCache "github.com/amirhossein-karimi/todo/internal/cache"
	"github.com/amirhossein-karimi/todo/internal/errors"
	"github.com/amirhossein-karimi/todo/internal/mocks"
	"github.com/amirhossein-karimi/todo/internal/models"
	"github.com/amirhossein-karimi/todo/internal/requests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

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
				Assignee:    "John Doe",
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
				Assignee:    "John Doe",
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
			repo := new(mocks.MockTodoRepository)
			cache := new(mocks.MockCache)
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

				repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
				cache.AssertNotCalled(t, "Delete", mock.Anything)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.createResult, result)

			repo.AssertExpectations(t)
		})
	}
}

func TestTodoService_Delete(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		uuid        uuid.UUID
		deleteErr   error
		cacheErr    error
		expectedErr error
	}{
		{
			name:        "should delete todo successfully",
			uuid:        uuid.New(),
			deleteErr:   nil,
			cacheErr:    nil,
			expectedErr: nil,
		},
		{
			name:        "should return error when todo not found",
			uuid:        uuid.New(),
			deleteErr:   gorm.ErrRecordNotFound,
			cacheErr:    nil,
			expectedErr: errors.ErrTodoNotFound,
		},
		{
			name:        "should return repository error",
			uuid:        uuid.New(),
			deleteErr:   gorm.ErrInvalidDB,
			cacheErr:    nil,
			expectedErr: gorm.ErrInvalidDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.MockTodoRepository)
			cache := new(mocks.MockCache)

			service := NewTodoService(repo, cache)

			repo.
				On("Delete", ctx, tt.uuid).
				Return(tt.deleteErr).
				Once()

			if tt.deleteErr == nil {
				key := realCache.CreateKeyName(realCache.TodosKey, "*")

				cache.On("Delete", ctx, key).Return(tt.cacheErr).Once()
			}

			err := service.Delete(ctx, tt.uuid)

			assert.ErrorIs(t, err, tt.expectedErr)

			repo.AssertExpectations(t)
			cache.AssertExpectations(t)
		})
	}
}

func TestTodoService_Update(t *testing.T) {
	ctx := context.Background()

	title := "Updated title"
	description := "Updated description"
	assignee := "John"
	priority := 2
	status := 1

	tests := []struct {
		name string
		uuid uuid.UUID
		req  *requests.TodoUpdateRequest

		findTodo    *models.Todo
		findTodoErr error

		findTitle    *models.Todo
		findTitleErr error

		updateResult *models.Todo
		updateErr    error

		expectedErr error
	}{
		{
			name: "should update todo successfully",
			uuid: uuid.New(),

			req: &requests.TodoUpdateRequest{
				Title:       &title,
				Description: &description,
				Assignee:    &assignee,
				Priority:    &priority,
				Status:      &status,
			},

			findTodo: &models.Todo{
				Title: "Old title",
			},

			findTitleErr: gorm.ErrRecordNotFound,

			updateResult: &models.Todo{
				Title: title,
			},

			updateErr: nil,
		},

		{
			name: "should return todo not found",
			uuid: uuid.New(),

			req: &requests.TodoUpdateRequest{
				Title: &title,
			},

			findTodoErr: gorm.ErrRecordNotFound,

			expectedErr: errors.ErrTodoNotFound,
		},

		{
			name: "should return error when title exists",
			uuid: uuid.New(),

			req: &requests.TodoUpdateRequest{
				Title: &title,
			},

			findTodo: &models.Todo{
				Title: "Old title",
			},

			findTitle: &models.Todo{
				Title: title,
			},

			findTitleErr: nil,

			expectedErr: errors.ErrTodoWithThisTitleAlreadyExists,
		},

		{
			name: "should return update error",
			uuid: uuid.New(),

			req: &requests.TodoUpdateRequest{
				Title: &title,
			},

			findTodo: &models.Todo{
				Title: "Old title",
			},

			findTitleErr: gorm.ErrRecordNotFound,

			updateErr: gorm.ErrInvalidDB,

			expectedErr: gorm.ErrInvalidDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			repo := new(mocks.MockTodoRepository)
			cache := new(mocks.MockCache)

			service := NewTodoService(repo, cache)

			repo.On("FindByUUID", ctx, tt.uuid).Return(tt.findTodo, tt.findTodoErr)

			if tt.findTodoErr == nil {

				repo.On("FindByTitleExcludingUUID", ctx, tt.uuid, *tt.req.Title).Return(tt.findTitle, tt.findTitleErr)

				if tt.findTitle == nil {

					repo.On("UpdateByUUID", ctx, tt.uuid, mock.Anything).Return(tt.updateResult, tt.updateErr)

					if tt.updateErr == nil {

						key := realCache.CreateKeyName(realCache.TodosKey, "*")

						cache.On("Delete", ctx, key).Return(nil)
					}
				}
			}

			result, err := service.Update(ctx, tt.uuid, tt.req)

			assert.ErrorIs(t, err, tt.expectedErr)

			if tt.expectedErr == nil {
				assert.Equal(t, tt.updateResult, result)
			}

			repo.AssertExpectations(t)
			cache.AssertExpectations(t)
		})
	}
}

func TestTodoService_List(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		page          int
		status        string
		priority      string
		assignee      string
		cacheTodosErr error
		cacheTotalErr error
		todos         []*models.Todo
		total         int64
		repoErr       error
		expectedTodos []*models.Todo
		expectedTotal int64
		expectedErr   error
	}{
		{
			name:     "should return todos from cache",
			page:     1,
			status:   "1",
			priority: "1",
			assignee: "john",
			expectedTodos: []*models.Todo{
				{
					Title: "Test todo",
				},
			},
			expectedTotal: 1,
		},
		{
			name:          "should get todos from repository when cache miss",
			page:          1,
			status:        "1",
			cacheTodosErr: gorm.ErrRecordNotFound,
			todos: []*models.Todo{
				{
					Title: "Learn Go",
				},
			},
			total: 1,
			expectedTodos: []*models.Todo{
				{
					Title: "Learn Go",
				},
			},
			expectedTotal: 1,
		},
		{
			name:          "should return repository error",
			page:          1,
			cacheTodosErr: gorm.ErrRecordNotFound,
			repoErr:       gorm.ErrInvalidDB,
			expectedErr:   gorm.ErrInvalidDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			repo := new(mocks.MockTodoRepository)
			cache := new(mocks.MockCache)

			service := NewTodoService(repo, cache)

			key := realCache.CreateKeyName(realCache.TodosKey, "page", tt.page, "status", tt.status, "priority", tt.priority, "assignee", tt.assignee)

			totalKey := realCache.CreateKeyName(realCache.TodosKey, "total")

			if tt.cacheTodosErr == nil {

				cache.On("Get", ctx, key, mock.Anything).Return(nil, tt.expectedTodos)

				cache.On("Get", ctx, totalKey, mock.Anything).Return(nil, tt.expectedTotal)

			} else {

				cache.On("Get", ctx, key, mock.Anything).Return(tt.cacheTodosErr, nil)

				repo.
					On("List", ctx, tt.page, tt.status, tt.priority, tt.assignee).Return(tt.todos, tt.total, tt.repoErr)

				if tt.repoErr == nil && len(tt.todos) > 0 {

					cache.On("Set", ctx, key, tt.todos, 10*time.Minute).Return(nil)

					cache.On("Set", ctx, totalKey, tt.total, 10*time.Minute).Return(nil)
				}
			}

			result, total, err := service.List(ctx, tt.page, tt.status, tt.priority, tt.assignee)

			assert.ErrorIs(t, err, tt.expectedErr)

			if tt.expectedErr == nil {
				assert.Equal(t, tt.expectedTodos, result)
				assert.Equal(t, tt.expectedTotal, total)
			}

			repo.AssertExpectations(t)
			cache.AssertExpectations(t)
		})
	}
}
