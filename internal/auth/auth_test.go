package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestVerifyAccessToken(t *testing.T) {
	correctKey := []byte("secret-key")
	wrongKey := []byte("wrong-key")

	tests := []struct {
		name            string
		signingKey      []byte
		verificationKey []byte
		issuer          string
		expectedExpiry  time.Duration
		wantErr         bool
	}{
		{
			name:            "valid token",
			signingKey:      correctKey,
			verificationKey: correctKey,
			issuer:          JWT_ISSUER,
			expectedExpiry:  5 * time.Minute,
			wantErr:         false,
		},
		{
			name:            "wrong verification key",
			signingKey:      correctKey,
			verificationKey: wrongKey,
			issuer:          JWT_ISSUER,
			expectedExpiry:  5 * time.Minute,
			wantErr:         true,
		},
		{
			name:            "token signed with wrong key",
			signingKey:      wrongKey,
			verificationKey: correctKey,
			issuer:          JWT_ISSUER,
			expectedExpiry:  5 * time.Minute,
			wantErr:         true,
		},
		{
			name:            "wrong issuer",
			signingKey:      correctKey,
			verificationKey: correctKey,
			issuer:          "wrong-issuer",
			expectedExpiry:  5 * time.Minute,
			wantErr:         true,
		},
		{
			name:            "expired token",
			signingKey:      correctKey,
			verificationKey: correctKey,
			issuer:          JWT_ISSUER,
			expectedExpiry:  -1 * time.Hour,
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := makeUser()

			token, err := makeJWT(generateJWTParams{
				issuer: tt.issuer,
				exp:    time.Now().Add(tt.expectedExpiry),
				key:    tt.signingKey,
			}, user)
			require.NoError(t, err)

			claims, err := VerifyAccessToken(token, tt.verificationKey)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, user.ID, claims.UserID)
		})
	}
}

func TestComparePassword(t *testing.T) {
	tests := []struct {
		name              string
		password          string
		passwordToCompare string
		isValid           bool
		wantErr           bool
	}{
		{
			name:              "valid password",
			password:          "super-secure-123",
			passwordToCompare: "super-secure-123",
			isValid:           true,
			wantErr:           false,
		},
		{
			name:              "invalid password",
			password:          "super-secure-123",
			passwordToCompare: "wrong-password",
			isValid:           false,
			wantErr:           false,
		},
		{
			name:              "empty password string matching",
			password:          "",
			passwordToCompare: "",
			isValid:           true,
			wantErr:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			require.NoError(t, err, "Hashing the password should not fail")

			valid, err := ComparePassword(ComparePasswordParams{
				Password:   tt.passwordToCompare,
				StoredHash: hash,
			})

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			require.Equal(t, tt.isValid, valid)
		})
	}
}
