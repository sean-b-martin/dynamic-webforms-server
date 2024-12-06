package service

import (
	"github.com/google/uuid"
	"github.com/sean-b-martin/dynamic-webforms-server/model"
	"github.com/uptrace/bun"
)

type SchemaService interface {
	GetSchemas(formID uuid.UUID) ([]model.FormSchemaModel, error)
}

func NewSchemaService(db *bun.DB) SchemaService {
	return &SchemaServiceImpl{db: db}
}

type SchemaServiceImpl struct {
	db        *bun.DB
	dbService GenericDBService[model.FormSchemaModel]
}

func (s *SchemaServiceImpl) GetSchemas(formID uuid.UUID) ([]model.FormSchemaModel, error) {
	return s.dbService.GetModels("form_id = ?", formID)
}
