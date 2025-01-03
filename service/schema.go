package service

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/sean-b-martin/dynamic-webforms-server/database"
	"github.com/sean-b-martin/dynamic-webforms-server/model"
	"github.com/uptrace/bun"
)

type SchemaService interface {
	GetSchemas(formID uuid.UUID) ([]model.FormSchemaModel, error)
	GetSchema(formID uuid.UUID, schemaID uuid.UUID) (model.FormSchemaModel, error)
	CreateSchema(formID uuid.UUID, userID uuid.UUID, formSchema model.FormSchemaModel) error
	UpdateSchema(username uuid.UUID, formID uuid.UUID, schemaID uuid.UUID, schemaData map[string]interface{}) error
	DeleteSchema(userID uuid.UUID, formID uuid.UUID, schemaID uuid.UUID) error
}

func NewSchemaService(db *bun.DB) SchemaService {
	return &SchemaServiceImpl{db: db, dbService: NewGenericDBService[model.FormSchemaModel](db)}
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

func (s *SchemaServiceImpl) CreateSchema(userID uuid.UUID, formID uuid.UUID, schema model.FormSchemaModel) error {
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer database.TXLogErrRollback(&tx)

	schema.FormID = formID
	if err := isFormOwner(&tx, formID, userID); err != nil {
		return err
	}

	_, err = tx.NewInsert().Model(&schema).Column("title", "version", "schema", "read_only", "form_id").
		Exec(context.Background())
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *SchemaServiceImpl) UpdateSchema(username uuid.UUID, formID uuid.UUID, schemaID uuid.UUID, schemaData map[string]interface{}) error {
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer database.TXLogErrRollback(&tx)

	if err := isFormOwner(&tx, formID, username); err != nil {
		return err
	}

	query := tx.NewUpdate().Model((*model.FormSchemaModel)(nil)).Where("id = ? AND form_id = ? ", schemaID, formID)
	for k, v := range schemaData {
		query.SetColumn(k, "?", v)
	}
	if _, err := query.Exec(context.Background()); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *SchemaServiceImpl) DeleteSchema(userID uuid.UUID, formID uuid.UUID, schemaID uuid.UUID) error {
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer database.TXLogErrRollback(&tx)

	if err := isFormOwner(&tx, formID, userID); err != nil {
		return err
	}

	res, err := tx.NewDelete().Model((*model.FormSchemaModel)(nil)).
		Where("id = ? AND form_id = ? ", schemaID, formID).Exec(context.Background())

	if err != nil {
		return err
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}

	return tx.Commit()
}

func isFormOwner(tx *bun.Tx, formID uuid.UUID, userID uuid.UUID) error {
	var form model.FormModel
	err := tx.NewSelect().Model((*model.FormModel)(nil)).Column("user_id").Where("id = ?", formID).
		Scan(context.Background(), &form)

	if err != nil {
		return err
	}

	if form.UserID != userID.String() {
		return ErrNoPermission
	}

	return nil
}
