package utils

import "github.com/jackc/pgx/v5/pgtype"

func NullStringToString(ns pgtype.Text) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
