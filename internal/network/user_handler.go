package network

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"social-media-backend/internal/adapters/storage"
	"social-media-backend/internal/apperrors"
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
		_ = c.Error(apperrors.MissingUserClaimsInContext)
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

// GetSelfUser godoc
//
//	@Summary		Gets the logged in user
//	@Description	Allows an authenticated user to Get their data using the jwt
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			Authorization	header		string			true	"Type 'Bearer _JWT_'"
//
//	@Success		200		{object}	map[string]any		"Password updated successfully"
//	@Failure		400		{object}	Response
//	@Failure		401		{object}	Response
//	@Failure		500		{object}	Response
//	@Router			/users/me/ [get]
func (h *UserHandler) GetSelfUser(c *gin.Context) {
	userData, ok := getClaims(c)
	if !ok {
		_ = c.Error(apperrors.MissingUserClaimsInContext)
		return
	}

	user, err := h.svc.GetUserByID(c.Request.Context(), userData.UserID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	OK(c, user)
}

// UpdateSelfAvatar godoc
//
//	@Summary		Updates the user avatar
//	@Description	Allows an authenticated user to upload and update their profile picture
//	@Tags			users
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		BearerAuth
//
//	@Param			Authorization	header		string	true	"Bearer JWT token"
//	@Param			file			formData	file	true	"Avatar image file"
//
//	@Success		200		{object}	map[string]any	"Avatar updated successfully"
//	@Failure		400		{object}	Response
//	@Failure		401		{object}	Response
//	@Failure		500		{object}	Response
//	@Router			/users/me/avatar [put]
func (h *UserHandler) UpdateSelfAvatar(c *gin.Context) {
	userData, ok := getClaims(c)
	if !ok {
		_ = c.Error(apperrors.MissingUserClaimsInContext)
		return
	}

	imageFile, err := c.FormFile("image")
	if err != nil {
		_ = c.Error(apperrors.MissingFile)
		return
	}

	contentType, err := ValidateFile(
		imageFile,
		5<<20,
		[]string{
			storage.ContentTypeJPEG.String(),
			storage.ContentTypePNG.String(),
		},
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	file, err := imageFile.Open()
	if err != nil {
		_ = c.Error(err)
		return
	}
	defer file.Close()

	var img image.Image

	switch contentType {
	case storage.ContentTypeJPEG.String():
		img, err = jpeg.Decode(file)
	case storage.ContentTypePNG.String():
		img, err = png.Decode(file)
	default:
		_ = c.Error(apperrors.InvalidMime)
		return
	}

	if err != nil {
		_ = c.Error(err)
		return
	}

	var buf bytes.Buffer

	err = jpeg.Encode(&buf, img, &jpeg.Options{
		Quality: 85,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	imgUrl, err := h.svc.UpdateUserAvatar(c.Request.Context(), userData.UserID, &buf)
	if err != nil {
		_ = c.Error(err)
		return
	}

	OK(c, gin.H{
		"image_url": imgUrl,
	})
}
