package utils

import "github.com/jackc/pgx/v5/pgtype"

func NullStringToString(ns pgtype.Text) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func PgBigIntToInt64(v pgtype.Int8) int64 {
	if v.Valid {
		return v.Int64
	}
	return 0
}
