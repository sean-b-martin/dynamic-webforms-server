package database

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/uptrace/bun"
)

func TXLogErrRollback(tx *bun.Tx) {
	if err := tx.Rollback(); err != nil {
		log.Error(err)
	}
}
