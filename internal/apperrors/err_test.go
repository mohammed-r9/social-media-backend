package apperrors

import (
	"testing"
)

func TestAuthErrTable(t *testing.T) {
	if len(authErrorTable) != int(authErrorCount) {
		t.Fatalf(
			"auth error table size mismatch: got %d, expected %d",
			len(authErrorTable),
			authErrorCount,
		)
	}
	for e := range int(authErrorCount) {
		info := authErrorTable[e]

		if info.code == "" ||
			info.message == "" ||
			info.status == 0 {
			t.Fatalf(
				"invalid entry at %v (%v): %+v",
				e,
				authError(e),
				info,
			)
		}
	}
}

func TestUserErrTable(t *testing.T) {
	if len(userErrorTable) != int(userErrorCount) {
		t.Fatalf(
			"user error table size mismatch: got %d, expected %d",
			len(userErrorTable),
			userErrorCount,
		)
	}

	for e := range int(userErrorCount) {
		info := userErrorTable[e]

		if info.code == "" ||
			info.message == "" ||
			info.status == 0 {
			t.Fatalf(
				"invalid entry at %v (%v): %+v",
				e,
				userError(e),
				info,
			)
		}
	}
}

func TestSessionErrTable(t *testing.T) {
	if len(sessionErrorTable) != int(sessionErrorCount) {
		t.Fatalf(
			"session error table size mismatch: got %d, expected %d",
			len(sessionErrorTable),
			sessionErrorCount,
		)
	}

	for e := range int(sessionErrorCount) {
		info := sessionErrorTable[e]

		if info.code == "" ||
			info.message == "" ||
			info.status == 0 {
			t.Fatalf(
				"invalid entry at %v (%v): %+v",
				e,
				sessionError(e),
				info,
			)
		}
	}
}

func TestNetworkErrTable(t *testing.T) {
	if len(networkErrorTable) != int(networkErrorCount) {
		t.Fatalf(
			"network error table size mismatch: got %d, expected %d",
			len(networkErrorTable),
			networkErrorCount,
		)
	}

	for e := range int(networkErrorCount) {
		info := networkErrorTable[e]

		if info.code == "" ||
			info.message == "" ||
			info.status == 0 {
			t.Fatalf(
				"invalid entry at %v (%v): %+v",
				e,
				networkError(e),
				info,
			)
		}
	}
}

func TestEnvErrTable(t *testing.T) {
	if len(envErrorTable) != int(envErrorCount) {
		t.Fatalf(
			"env error table size mismatch: got %d, expected %d",
			len(envErrorTable),
			envErrorCount,
		)
	}

	for e := range int(envErrorCount) {
		info := envErrorTable[e]

		if info.code == "" ||
			info.message == "" ||
			info.status == 0 {
			t.Fatalf(
				"invalid entry at %v (%v): %+v",
				e,
				envError(e),
				info,
			)
		}
	}
}

func TestDatabaseErrTable(t *testing.T) {
	if len(databaseErrorTable) != int(dbErrorCount) {
		t.Fatalf(
			"database error table size mismatch: got %d, expected %d",
			len(databaseErrorTable),
			dbErrorCount,
		)
	}

	for e := range int(dbErrorCount) {
		info := databaseErrorTable[e]

		if info.code == "" ||
			info.message == "" ||
			info.status == 0 {
			t.Fatalf(
				"invalid entry at %v (%v): %+v",
				e,
				databaseError(e),
				info,
			)
		}
	}
}
