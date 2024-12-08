package service

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type GenericDBService[T any] interface {
	GetModelByID(id uuid.UUID) (T, error)
	GetModel(whereQuery string, args ...interface{}) (T, error)
	GetModels(whereQuery string, args ...interface{}) ([]T, error)
	InsertModel(model T, columns ...string) error
	UpdateModel(model T, id uuid.UUID, columns ...string) error
	DeleteModelByID(id uuid.UUID) error
}

type genericDBServiceImpl[T any] struct {
	db *bun.DB
}

func NewGenericDBService[T any](db *bun.DB) GenericDBService[T] {
	return &genericDBServiceImpl[T]{db: db}
}

func (g *genericDBServiceImpl[T]) GetModels(whereQuery string, args ...interface{}) ([]T, error) {
	var models []T

	query := g.db.NewSelect().Model((*T)(nil))
	if whereQuery != "" {
		query.Where(whereQuery, args...)
	}

	err := query.Scan(context.Background(), &models)
	return models, err
}

func (g *genericDBServiceImpl[T]) GetModelByID(id uuid.UUID) (T, error) {
	var model T
	err := g.db.NewSelect().Model(&model).Where("id = ?", id).Scan(context.Background())
	return model, err
}

func (g *genericDBServiceImpl[T]) GetModel(whereQuery string, args ...interface{}) (T, error) {
	var model T
	err := g.db.NewSelect().Model(&model).Where(whereQuery, args...).Scan(context.Background())
	return model, err
}

func (g *genericDBServiceImpl[T]) InsertModel(model T, columns ...string) error {
	_, err := g.db.NewInsert().Model(&model).Column(columns...).Exec(context.Background())
	return err
}

func (g *genericDBServiceImpl[T]) UpdateModel(model T, id uuid.UUID, columns ...string) error {
	if res, err := g.db.NewUpdate().Model(&model).Column(columns...).Where("id = ?", id).
		Exec(context.Background()); err != nil {
		return err
	} else if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (g *genericDBServiceImpl[T]) DeleteModelByID(id uuid.UUID) error {
	var model T
	if res, err := g.db.NewDelete().Model(&model).Where("id = ?", id).Exec(context.Background()); err != nil {
		return err
	} else if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
