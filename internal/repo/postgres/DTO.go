package postgres

import (
	"social-media-backend/internal/domain"
	"time"
)

type UserSessionDTO struct {
	Session    domain.Session
	VerifiedAt *time.Time
}
