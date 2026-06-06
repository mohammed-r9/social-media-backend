package auth

// HashTokens requires both RefreshToken and CsrfToken to be set before being called,
func (t *SessionTokens) ToHash() SessionTokenHashes {
	return SessionTokenHashes{
		RefreshHash: HashToken(t.RefreshToken),
		CsrfHash:    HashToken(t.CsrfToken),
	}
}
