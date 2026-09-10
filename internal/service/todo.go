package service

import (
	"context"

	"time"

	"github.com/amirhossein-karimi/todo/internal/cache"
	"github.com/amirhossein-karimi/todo/internal/errors"
	"github.com/amirhossein-karimi/todo/internal/models"
	"github.com/amirhossein-karimi/todo/internal/repository"
	"github.com/amirhossein-karimi/todo/internal/requests"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TodoService interface {
	Create(ctx context.Context, todo *requests.TodoRequest) (uuid.UUID, error)
	List(ctx context.Context, page int) ([]*models.Todo, int64, error)
}

type todoService struct {
	repo  repository.TodoRepository
	cache cache.Cache
}

func NewTodoService(repo repository.TodoRepository, cache cache.Cache) TodoService {
	return &todoService{
		repo:  repo,
		cache: cache,
	}
}

func (s *todoService) Create(ctx context.Context, todo *requests.TodoRequest) (uuid.UUID, error) {

	todoModel, err := s.repo.FindByTitle(ctx, todo.Title)
	if err != nil && err != gorm.ErrRecordNotFound {
		return uuid.Nil, err
	}
	if todoModel != nil {
		return uuid.Nil, errors.ErrTodoWithThisTitleAlreadyExists
	}

	uuidData := uuid.New()
	now := time.Now()

	uuidField, err := s.repo.Create(ctx, &models.Todo{
		Title:       todo.Title,
		Description: todo.Description,
		Priority:    todo.Priority,
		Status:      models.TodoStatusTodo,
		UUID:        uuidData,
		ToDoStartAt: &now,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		return uuid.Nil, err
	}

	key := cache.CreateKeyName(cache.TodosKey, "*")

	err = s.cache.Delete(ctx, key)
	if err != nil {
		return uuid.Nil, err
	}
	return uuidField, nil
}

func (s *todoService) List(ctx context.Context, page int) ([]*models.Todo, int64, error) {
	const ttl = 10 * time.Minute

	key := cache.CreateKeyName(cache.TodosKey, "page", page)
	totalKey := cache.CreateKeyName(cache.TodosKey, "total")

	var todos []*models.Todo
	var total int64

	if err := s.cache.Get(ctx, key, &todos); err == nil {
		if err := s.cache.Get(ctx, totalKey, &total); err == nil {
			return todos, total, nil
		}
	}

	todos, total, err := s.repo.List(ctx, page)
	if err != nil {
		return nil, 0, err
	}

	_ = s.cache.Set(ctx, key, todos, ttl)
	_ = s.cache.Set(ctx, totalKey, total, ttl)

	return todos, total, nil
}
