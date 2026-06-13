package repo

import "github.com/jackc/pgx/v5/pgtype"

func textToString(ns pgtype.Text) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func int8ToInt64(v pgtype.Int8) int64 {
	if v.Valid {
		return v.Int64
	}
	return 0
}

func stringPtrToTex(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{
		String: *s,
		Valid:  true,
	}
}

func stringToTex(s string) pgtype.Text {
	return pgtype.Text{
		String: s,
		Valid:  true,
	}
}
