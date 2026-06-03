package network

import (
	"social-media-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

type updatePasswordReq struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// UpdateSelfPassword godoc
//
//	@Summary		Update user password
//	@Description	Allows an authenticated user to change their password by providing the old password and a new password.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			Authorization	header		string			true	"Type 'Bearer _JWT_'"
//
//	@Param			request	body		updatePasswordReq	true	"Password update payload"
//	@Success		200		{object}	map[string]any		"Password updated successfully"
//	@Failure		400		{object}	Response
//	@Failure		401		{object}	Response
//	@Failure		500		{object}	Response
//	@Router			/users/me/password [put]
func (h *UserHandler) UpdateSelfPassword(c *gin.Context) {
	userData, ok := getClaims(c)
	if !ok {
		_ = c.Error(errMissingUserClaimsInContext)
		return
	}

	var req updatePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	err := h.svc.UpdateUserPassword(c.Request.Context(), service.UpdateUserPasswordParams{
		UserID:      userData.UserID,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})

	if err != nil {
		_ = c.Error(err)
		return
	}

	OK(c, gin.H{
		"message": "Password updated successfully",
	})
}
