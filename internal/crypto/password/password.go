package password

import (
	"github.com/alexedwards/argon2id"
)

var DefaultParams = &argon2id.Params{
	Memory:      64 * 1024,
	Iterations:  3,
	Parallelism: 4,
	KeyLength:   32,
	SaltLength:  16,
}

func HashPassword(plainText string) (string, error) {
	return argon2id.CreateHash(plainText, DefaultParams)
}

type ComparePasswordParams struct {
	Password   string
	StoredHash string
}

func ComparePassword(params ComparePasswordParams) (bool, error) {
	return argon2id.ComparePasswordAndHash(params.Password, params.StoredHash)
}
