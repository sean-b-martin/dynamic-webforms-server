package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/sean-b-martin/dynamic-webforms-server/database"
	"github.com/sean-b-martin/dynamic-webforms-server/model"
	"github.com/uptrace/bun"
)

type SchemaService interface {
	GetSchemas(formID uuid.UUID) ([]model.FormSchemaModel, error)
	GetSchema(formID uuid.UUID, schemaID uuid.UUID) (model.FormSchemaModel, error)
	CreateSchema(formID uuid.UUID, userID uuid.UUID, formSchema model.FormSchemaModel) error
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

func (s *SchemaServiceImpl) GetSchema(formID uuid.UUID, schemaID uuid.UUID) (model.FormSchemaModel, error) {
	return s.dbService.GetModel("id = ? AND form_id = ? ", schemaID, formID)
}

func (s *SchemaServiceImpl) CreateSchema(formID uuid.UUID, userID uuid.UUID, schema model.FormSchemaModel) error {
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer database.TXLogErrRollback(&tx)

	schema.FormID = formID
	if !isFormOwner(&tx, formID, userID) {
		return ErrNoPermission
	}

	_, err = tx.NewInsert().Model(&schema).Column("title", "version", "schema", "read_only", "form_id").
		Exec(context.Background())
	if err != nil {
		return err
	}

	return tx.Commit()
}

func isFormOwner(tx *bun.Tx, formID uuid.UUID, userID uuid.UUID) bool {
	var form model.FormModel
	err := tx.NewSelect().Model((*model.FormModel)(nil)).Column("user_id").Where("id = ?", formID).
		Scan(context.Background(), form)

	if err != nil {
		return false
	}

	if form.UserID != userID.String() {
		return false
	}

	return true
}
