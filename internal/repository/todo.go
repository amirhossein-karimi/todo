package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/amirhossein-karimi/todo/internal/models"
	"gorm.io/gorm"
)

type TodoRepository interface {
	Create(ctx context.Context, todo *models.Todo) (uuid.UUID, error)
	FindByTitle(ctx context.Context, title string) (*models.Todo, error)
	FindByUUID(ctx context.Context, uuid uuid.UUID) (*models.Todo, error)
	List(ctx context.Context, page int, status string, priority string, assignee string) ([]*models.Todo, int64, error)
	Delete(ctx context.Context, uuid uuid.UUID) error
}

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{
		db: db,
	}
}



func (r *todoRepository) FindByUUID(ctx context.Context, uuid uuid.UUID) (*models.Todo, error) {
	var todo models.Todo
	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&todo).Error
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *todoRepository) Delete(ctx context.Context, uuid uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(&models.Todo{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *todoRepository) List(ctx context.Context, page int, status string, priority string, assignee string) ([]*models.Todo, int64, error) {

	const pageSize = 10

	var todos []*models.Todo
	var total int64

	db := r.db.WithContext(ctx).Model(&models.Todo{})

	if status != "" {
		db = db.Where("status = ?", status)
	}
	if priority != "" {
		db = db.Where("priority = ?", priority)
	}
	if assignee != "" {
		db = db.Where("assignee = ?", assignee)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	if err := db.
		Order("id DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&todos).Error; err != nil {
		return nil, 0, err
	}

	return todos, total, nil
}

func (r *todoRepository) Create(ctx context.Context, todo *models.Todo) (uuid.UUID, error) {
	err := r.db.WithContext(ctx).Create(todo).Error
	if err != nil {
		return uuid.Nil, err
	}
	return todo.UUID, nil
}

func (r *todoRepository) FindByTitle(ctx context.Context, title string) (*models.Todo, error) {
	var todo models.Todo
	err := r.db.WithContext(ctx).Where("title = ?", title).First(&todo).Error
	if err != nil {
		return nil, err
	}
	return &todo, nil
}
