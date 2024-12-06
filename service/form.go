package service

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/sean-b-martin/dynamic-webforms-server/database"
	"github.com/sean-b-martin/dynamic-webforms-server/model"
	"github.com/uptrace/bun"
)

type FormService interface {
	GetFormsOfUser(userID uuid.UUID) ([]model.FormModel, error)
	GetForm(formID uuid.UUID) (model.FormModel, error)
	GetForms() ([]model.FormModel, error)
	CreateForm(userID uuid.UUID, title string) error
	UpdateForm(userID uuid.UUID, id uuid.UUID, title string) error
	DeleteForm(userID uuid.UUID, id uuid.UUID) error
}

type formServiceImpl struct {
	db        *bun.DB
	dbService GenericDBService[model.FormModel]
}

func NewFormService(db *bun.DB) FormService {
	return &formServiceImpl{db: db, dbService: NewGenericDBService[model.FormModel](db)}
}

func (f *formServiceImpl) GetFormsOfUser(userID uuid.UUID) ([]model.FormModel, error) {
	var forms []model.FormModel

	err := f.db.NewSelect().Model((*model.FormModel)(nil)).Where("user_id = ?", userID).
		Scan(context.Background(), &forms)
	if err != nil {
		return nil, err
	}

	return forms, nil
}

func (f *formServiceImpl) GetForms() ([]model.FormModel, error) {
	return f.dbService.GetModels("")
}

func (f *formServiceImpl) GetForm(formID uuid.UUID) (model.FormModel, error) {
	return f.dbService.GetModelByID(formID)
}

func (f *formServiceImpl) CreateForm(userID uuid.UUID, title string) error {
	form := model.FormModel{Title: title, UserID: userID.String()}
	return f.dbService.InsertModel(form, "user_id", "title")
}

func (f *formServiceImpl) UpdateForm(userID uuid.UUID, id uuid.UUID, title string) error {
	var form model.FormModel

	tx, err := f.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	defer database.TXLogErrRollback(&tx)

	err = tx.NewSelect().Model(&form).Where("id = ?", id).Scan(context.Background())
	if err != nil {
		return err
	}

	if form.UserID != userID.String() {
		return ErrNoPermission
	}

	form.Title = title
	_, err = tx.NewUpdate().Model(&form).Column("title").Where("id = ?", id).Exec(context.Background())

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (f *formServiceImpl) DeleteForm(userID uuid.UUID, id uuid.UUID) error {
	var form model.FormModel
	tx, err := f.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer database.TXLogErrRollback(&tx)

	if err := tx.NewSelect().Model(&form).Where("id = ?", id).Scan(context.Background()); err != nil {
		return err
	}

	if form.UserID != userID.String() {
		return ErrNoPermission
	}

	if res, err := tx.NewDelete().Model((*model.FormModel)(nil)).Where("id = ?", id).Exec(context.Background()); err != nil {
		return err
	} else if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}

	return tx.Commit()
}
