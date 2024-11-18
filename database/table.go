package database

import (
	"context"
	"fmt"
	"github.com/uptrace/bun"
	"log"
)

func CreateTables(db *bun.DB) {
	// activate uuid extension to use uuid_generate_v4 db function
	if _, err := db.NewRaw(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).
		Exec(context.Background()); err != nil {
		log.Fatal(err)
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*UserModel)(nil)).Exec(context.Background()); err != nil {
		log.Fatal(fmt.Errorf("failed creating table for UserModel: %w", err))
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*FormModel)(nil)).
		ForeignKey(`("user_id") REFERENCES "users" ("id") ON DELETE SET NULL`).
		Exec(context.Background()); err != nil {
		log.Fatal(fmt.Errorf("failed creating table for FormModel: %w", err))
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*FormSchemaModel)(nil)).
		ForeignKey(`("form_id") REFERENCES "forms" ("id") ON DELETE CASCADE`).
		Exec(context.Background()); err != nil {
		log.Fatal(fmt.Errorf("failed creating table for FormSchemaModel: %w", err))
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*FormDataModel)(nil)).
		ForeignKey(`("user_id") REFERENCES "users" ("id") ON DELETE CASCADE`).
		ForeignKey(`("form_schema_id") REFERENCES "form_schemas" ("id") ON DELETE CASCADE`).
		Exec(context.Background()); err != nil {
		log.Fatal(fmt.Errorf("failed creating table for FormDataModel: %w", err))
	}

	if _, err := db.NewCreateTable().IfNotExists().Model((*FileMetadataModel)(nil)).
		ForeignKey(`("form_data_id") REFERENCES "form_data" ("id") ON DELETE CASCADE`).
		Exec(context.Background()); err != nil {
		log.Fatal(fmt.Errorf("failed creating table for FileMetadataModel: %w", err))
	}
}
