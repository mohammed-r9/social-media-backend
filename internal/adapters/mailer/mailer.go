package mailer

import (
	"social-media-backend/internal/env"

	"github.com/resend/resend-go/v3"
)

type ResendMailer struct {
	client *resend.Client
}

func NewResendMailer() *ResendMailer {
	return &ResendMailer{
		client: resend.NewClient(env.Config.RESEND_API_KEY),
	}
}

// LACKS IMPLEMENTATION
func (m *ResendMailer) Send(to, subject, body string) error {
	return nil
}
