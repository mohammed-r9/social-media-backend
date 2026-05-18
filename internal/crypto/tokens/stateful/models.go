package stateful

type SessionTokens struct {
	RefreshToken string
	CsrfToken    string
}

type SessionTokenHashes struct {
	RefreshHash string
	CsrfHash    string
}

// HashTokens requires both RefreshToken and CsrfToken to be set before being called,
func (t *SessionTokens) ToHash() SessionTokenHashes {
	return SessionTokenHashes{
		RefreshHash: hashToken(t.RefreshToken),
		CsrfHash:    hashToken(t.CsrfToken),
	}
}

func GenerateSessionTokens() SessionTokens {
	return SessionTokens{
		RefreshToken: generateOpaqueToken(32),
		CsrfToken:    generateOpaqueToken(32),
	}
}
