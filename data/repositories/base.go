package repositories

import (
	"github.com/Marcel-MD/clean-api/models"
	"gorm.io/gorm"
)

type BaseRepository[T any] interface {
	FindAll(query models.PaginationQuery) ([]T, error)
	FindById(id string) (T, error)
	Create(t *T) error
	Update(t *T) error
	Delete(t *T) error
}

func NewBaseRepository[T any](db *gorm.DB) BaseRepository[T] {
	return &baseRepository[T]{
		db: db,
	}
}

type baseRepository[T any] struct {
	db *gorm.DB
}

func (r *baseRepository[T]) FindAll(query models.PaginationQuery) ([]T, error) {
	var ts []T

	err := r.db.Scopes(paginate(query.Page, query.Size)).Find(&ts).Error
	return ts, err
}

func (r *baseRepository[T]) FindById(id string) (T, error) {
	var t T
	err := r.db.First(&t, "id = ?", id).Error
	return t, err
}

func (r *baseRepository[T]) Create(t *T) error {
	return r.db.Create(t).Error
}

func (r *baseRepository[T]) Update(t *T) error {
	return r.db.Save(t).Error
}

func (r *baseRepository[T]) Delete(t *T) error {
	return r.db.Delete(t).Error
}

func paginate(page int, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch {
		case size > 100:
			size = 100
		case size <= 0:
			size = 50
		}

		if page <= 0 {
			page = 1
		}

		offset := (page - 1) * size
		return db.Offset(offset).Limit(size)
	}
}
