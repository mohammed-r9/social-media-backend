package postgres

import (
	"social-media-backend/internal/domain"
	"time"
)

type UserSessionDTO struct {
	Session domain.Session
	User    struct {
		VerifiedAt *time.Time
	}
}
