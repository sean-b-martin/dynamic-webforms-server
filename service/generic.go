package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type GenericDBService[T any] interface {
	GetModelByID(id uuid.UUID) (T, error)
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

func (g *genericDBServiceImpl[T]) GetModelByID(id uuid.UUID) (T, error) {
	var model T
	err := g.db.NewSelect().Model(&model).Where("id = ?", id).Scan(context.Background())
	return model, err
}

func (g *genericDBServiceImpl[T]) InsertModel(model T, columns ...string) error {
	_, err := g.db.NewInsert().Model(&model).Column(columns...).Exec(context.Background())
	return err
}

func (g *genericDBServiceImpl[T]) UpdateModel(model T, id uuid.UUID, columns ...string) error {
	_, err := g.db.NewUpdate().Model(&model).Column(columns...).Where("id = ?", id).Exec(context.Background())
	return err
}

func (g *genericDBServiceImpl[T]) DeleteModelByID(id uuid.UUID) error {
	var model T
	_, err := g.db.NewDelete().Model(model).Where("id = ?", id).Exec(context.Background())
	return err
}