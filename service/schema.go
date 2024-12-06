package service

import "github.com/uptrace/bun"

type SchemaService interface {
}

func NewSchemaService(db *bun.DB) SchemaService {
	return &SchemaServiceImpl{db: db}
}

type SchemaServiceImpl struct {
	db *bun.DB
}
