package model

import (
	"database/sql"
	"encoding/json"
)

type NullString struct {
	sql.NullString
}

func ToNullString(s string) NullString {
	return NullString{sql.NullString{String: s, Valid: true}}
}

func (v *NullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	}
	return json.Marshal(nil)
}

type NullInt64 struct {
	sql.NullInt64
}

func ToNullInt64(i int) NullInt64 {
	return NullInt64{sql.NullInt64{Int64: int64(i), Valid: true}}
}

func (v *NullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	}
	return json.Marshal(nil)
}
