package service

import (
	"context"

	"time"

	"github.com/amirhossein-karimi/todo/internal/cache"
	"github.com/amirhossein-karimi/todo/internal/errors"
	"github.com/amirhossein-karimi/todo/internal/metrics"
	"github.com/amirhossein-karimi/todo/internal/models"
	"github.com/amirhossein-karimi/todo/internal/repository"
	"github.com/amirhossein-karimi/todo/internal/requests"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TodoService interface {
	Create(ctx context.Context, todo *requests.TodoRequest) (uuid.UUID, error)
	List(ctx context.Context, page int, status string, priority string, assignee string) ([]*models.Todo, int64, error)
	Delete(ctx context.Context, uuid uuid.UUID) error
	GetInfo(ctx context.Context, uuid uuid.UUID) (*models.Todo, error)
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

func (s *todoService) GetInfo(ctx context.Context, uuid uuid.UUID) (*models.Todo, error) {
	todo, err := s.repo.FindByUUID(ctx, uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrTodoNotFound
		}
		return nil, err
	}

	return todo, nil
}

func (s *todoService) Delete(ctx context.Context, uuid uuid.UUID) error {
	err := s.repo.Delete(ctx, uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrTodoNotFound
		}
		return err
	}

	key := cache.CreateKeyName(cache.TodosKey, "*")
	err = s.cache.Delete(ctx, key)
	if err != nil {
		return err
	}

	metrics.TasksCount.Dec()

	return nil
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
		Assignee:    todo.Assignee,
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
	metrics.TasksCount.Inc()
	return uuidField, nil
}

func (s *todoService) List(ctx context.Context, page int, status string, priority string, assignee string) ([]*models.Todo, int64, error) {
	const ttl = 10 * time.Minute

	key := cache.CreateKeyName(cache.TodosKey, "page", page, "status", status, "priority", priority, "assignee", assignee)
	totalKey := cache.CreateKeyName(cache.TodosKey, "total")

	var todos []*models.Todo
	var total int64

	if err := s.cache.Get(ctx, key, &todos); err == nil {
		if err := s.cache.Get(ctx, totalKey, &total); err == nil {
			return todos, total, nil
		}
	}

	todos, total, err := s.repo.List(ctx, page, status, priority, assignee)
	if err != nil {
		return nil, 0, err
	}
	if len(todos) > 0 {
		_ = s.cache.Set(ctx, key, todos, ttl)
		_ = s.cache.Set(ctx, totalKey, total, ttl)
	}

	return todos, total, nil
}
