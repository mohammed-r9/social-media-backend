package cache

import "social-media-backend/internal/repo"

func shouldCacheSession(s repo.UserSessionDTO) bool {
	return s.User.VerifiedAt != nil
}
