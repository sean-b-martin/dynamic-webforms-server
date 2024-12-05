package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/sean-b-martin/dynamic-webforms-server/auth"
	"github.com/sean-b-martin/dynamic-webforms-server/model"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"
)

type UserService interface {
	RegisterUser(user model.UserModel) error
	LoginUser(user model.UserModel) (string, error)
	GetUserById(id uuid.UUID) (model.UserModel, error)
}

type userServiceImpl struct {
	db              *bun.DB
	dbService       GenericDBService[model.UserModel]
	passwordService *auth.PasswordService
	jwtService      *auth.JWTService
}

func NewUserService(db *bun.DB, passwordService *auth.PasswordService, jwtService *auth.JWTService) UserService {
	return &userServiceImpl{
		db:              db,
		dbService:       NewGenericDBService[model.UserModel](db),
		passwordService: passwordService,
		jwtService:      jwtService,
	}
}

func (s *userServiceImpl) GetUserByUsername(username string) (model.UserModel, error) {
	var user model.UserModel
	err := s.db.NewSelect().Model(&user).Column("id", "username", "password").
		Where("username = ?", username).Scan(context.Background())
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *userServiceImpl) GetUserById(id uuid.UUID) (model.UserModel, error) {
	return s.dbService.GetModelByID(id)
}

func (s *userServiceImpl) RegisterUser(user model.UserModel) error {
	var err error

	user.Password, err = s.passwordService.HashPassword(user.Password)
	if err != nil {
		return err
	}

	if err = s.dbService.InsertModel(user, "username", "password"); err != nil {
		var pgErr pgdriver.Error
		if errors.As(err, &pgErr) {
			if pgErr.IntegrityViolation() {
				return errors.New("username already exists")
			}
		}
		return errors.New("error creating user")
	}
	return nil
}

func (s *userServiceImpl) LoginUser(user model.UserModel) (string, error) {
	dbUser, err := s.GetUserByUsername(user.Username)
	if err != nil {
		return "", err
	}

	if err := s.passwordService.VerifyPassword(dbUser.Password, user.Password); err != nil {
		return "", errors.New("invalid password")
	}

	return s.jwtService.NewToken(dbUser.ID)
}

// TODO Update

func (s *userServiceImpl) DeleteUser(id uuid.UUID) error {
	return s.dbService.DeleteModelByID(id)
}
