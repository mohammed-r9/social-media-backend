package storage

import (
	"context"
	"io"
)

type Storage interface {
	Upload(ctx context.Context, key string, body io.Reader, contentType ContentType) error
	Delete(ctx context.Context, key string) error
	GetURL(ctx context.Context, key string) (string, error)
}

type ContentType int

const (
	ContentTypePNG ContentType = iota
	ContentTypeJPEG
	ContentTypeGIF
	ContentTypeMP4
	ContentTypeWebM
)

func (c ContentType) String() string {
	switch c {
	case ContentTypePNG:
		return "image/png"
	case ContentTypeJPEG:
		return "image/jpeg"
	case ContentTypeGIF:
		return "image/gif"
	case ContentTypeMP4:
		return "video/mp4"
	case ContentTypeWebM:
		return "video/webm"

	default:
		return "application/octet-stream"
	}
}
