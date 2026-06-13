package domain

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserParams struct {
	ID           uuid.UUID
	Email        string
	Username     string
	Name         string
	PasswordHash string
}

type UpdatePasswordParams struct {
	ID           uuid.UUID
	PasswordHash string
}

type User struct {
	ID           uuid.UUID
	Email        string
	Username     string
	Phone        string
	PasswordHash string `json:"-"`
	IsSuspended  bool
	VerifiedAt   *time.Time
	Profile      Profile
}

type Profile struct {
	DisplayName    string
	Bio            string
	AvatarKey      string `json:"-"`
	Website        string
	FollowerCount  int64
	FollowingCount int64
	PostCount      int64
	AvatarURL      *string `json:"avatar_url"`
}

type UpdateProfileParams struct {
	UserID      uuid.UUID
	DisplayName *string
	Bio         *string
	Website     *string
	AvatarKey   *string // should implement it after adding s3
}

func (u *User) UpdateAvatarUrl(url string) {
	u.Profile.AvatarURL = &url
}
