package repository

import (
	"database/sql"
	"errors"
)

var NotFoundErr = errors.New("not found")

func sqlE(e error) error {
	if e == sql.ErrNoRows {
		return NotFoundErr
	}
	return e
}
