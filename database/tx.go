package database

import (
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2/log"
	"github.com/uptrace/bun"
)

func TXLogErrRollback(tx *bun.Tx) {
	if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
		log.Error(err)
	}
}
