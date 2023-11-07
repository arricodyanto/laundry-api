package utils

import (
	"database/sql"
)

func Validate(err error, message string, tx *sql.Tx) error {
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
