package database

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TableID struct {
	ID uuid.UUID `bun:"id,type:uuid,pk,default:uuid_generate_v4()"`
}

type UserModel struct {
	bun.BaseModel `bun:"table:users"`
	TableID
	Username string `bun:"username,type:varchar(128),notnull,unique"`
	Password string `bun:"password,type:varchar(60),notnull"`
}

type FormModel struct {
	bun.BaseModel `bun:"table:forms"`
	TableID
	UserID string `bun:"user_id,type:uuid"`
	Title  string `bun:"title,type:varchar(256),notnull"`
}

type FormSchemaModel struct {
	bun.BaseModel `bun:"table:form_schemas"`
	TableID
	FormID   uuid.UUID `bun:"form_id,type:uuid,notnull"`
	Title    string    `bun:"title,type:varchar(256),notnull"`
	Version  string    `bun:"version,type:varchar(64),notnull"`
	Schema   []byte    `bun:"schema,type:jsonb"`
	ReadOnly bool      `bun:"read_only,notnull,default:false"`
}

type FormDataModel struct {
	bun.BaseModel `bun:"table:form_data"`
	TableID
	UserID       uuid.UUID `bun:"user_id,type:uuid,notnull"`
	FormSchemaID uuid.UUID `bun:"form_schema_id,type:uuid,notnull"`
	Name         string    `bun:"name,type:varchar(64),notnull"`
	Data         []byte    `bun:"data,type:jsonb"`
}

type FileMetadataModel struct {
	bun.BaseModel `bun:"table:file_metadata"`
	TableID
	FormDataID         uuid.UUID `bun:"form_data_id,type:uuid,notnull"`
	OriginalFilename   string    `bun:"original_filename,type:varchar(256),notnull"`
	Path               string    `bun:"path,type:varchar(512),notnull"`
	MappingSchemaField int64     `bun:"mapping_schema_field,type:bigint,notnull"`
}
