package apperrors

import (
	"testing"
)

func TestAuthErrTable(t *testing.T) {
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
