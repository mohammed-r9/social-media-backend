package cache

// USER KEYS

func keyUserByID(ID string) string {
	return "user:" + ID
}

// SESSION KEYS

func keySessionByID(ID string) string {
	return "session:" + ID
}
