package mailer

type EmailService struct {
	mailer Mailer
}

func NewEmailService(mailer Mailer) *EmailService {
	return &EmailService{
		mailer: mailer,
	}
}

func (e *EmailService) SendVerificationEmail(to, token string) error {
	subject := "Verify your email address"
	body := buildVerificationEmailBody(token)

	return e.mailer.Send(to, subject, body)
}

func (e *EmailService) SendPasswordResetEmail(to, token string) error {
	subject := "Reset your password"
	body := buildPasswordResetEmailBody(token)

	return e.mailer.Send(to, subject, body)
}

func (e *EmailService) SendEmailChangeEmail(to, newEmail, token string) error {
	subject := "Confirm your new email address"
	body := buildEmailChangeEmailBody(newEmail, token)

	return e.mailer.Send(to, subject, body)
}

func (e *EmailService) SendSecurityAlertEmail(to, message string) error {
	subject := "Security alert"
	body := buildSecurityAlertBody(message)

	return e.mailer.Send(to, subject, body)
}
