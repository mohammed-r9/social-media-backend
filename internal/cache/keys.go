package cache

// USER KEYS

func keyUserByID(ID string) string {
	return "user:id" + ID
}

// SESSION KEYS

func keySessionByID(ID string) string {
	return "session:id" + ID
}
