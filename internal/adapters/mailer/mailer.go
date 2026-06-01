package mailer

import (
	"github.com/resend/resend-go/v3"
)

type ResendMailer struct {
	client *resend.Client
}

func NewResendMailer(apiKey string) *ResendMailer {
	return &ResendMailer{
		client: resend.NewClient(apiKey),
	}
}

// LACKS IMPLEMENTATION
func (m *ResendMailer) Send(to, subject, body string) error {
	return nil
}
