package model

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
	Username string `bun:"username,type:varchar(128),notnull,unique" json:"username"`
	Password string `bun:"password,type:varchar(60),notnull" json:"-"`
}

type FormModel struct {
	bun.BaseModel `bun:"table:forms"`
	TableID
	UserID string `bun:"user_id,type:uuid" json:"userID"`
	Title  string `bun:"title,type:varchar(256),notnull" json:"title"`
}

type FormSchemaModel struct {
	bun.BaseModel `bun:"table:form_schemas"`
	TableID
	FormID   uuid.UUID `bun:"form_id,type:uuid,notnull" json:"formID"`
	Title    string    `bun:"title,type:varchar(256),notnull" json:"title"`
	Version  string    `bun:"version,type:varchar(64),notnull" json:"version"`
	Schema   []byte    `bun:"schema,type:jsonb" json:"schema"`
	ReadOnly bool      `bun:"read_only,notnull,default:false" json:"readOnly"`
}

type FormDataModel struct {
	bun.BaseModel `bun:"table:form_data"`
	TableID
	UserID       uuid.UUID `bun:"user_id,type:uuid,notnull" json:"userID"`
	FormSchemaID uuid.UUID `bun:"form_schema_id,type:uuid,notnull" json:"formSchemaID"`
	Name         string    `bun:"name,type:varchar(64),notnull" json:"name"`
	Data         []byte    `bun:"data,type:jsonb" json:"data"`
}

type FileMetadataModel struct {
	bun.BaseModel `bun:"table:file_metadata"`
	TableID
	FormDataID         uuid.UUID `bun:"form_data_id,type:uuid,notnull" json:"formDataID"`
	OriginalFilename   string    `bun:"original_filename,type:varchar(256),notnull" json:"originalFilename"`
	Path               string    `bun:"path,type:varchar(512),notnull" json:"path"`
	MappingSchemaField int64     `bun:"mapping_schema_field,type:bigint,notnull" json:"mappingSchemaField"`
}
