package mail

import "context"

const applicationJsonContentType string = "application/json"

const apiUrlMailchimpSendEmail string = "https://mandrillapp.com/api/1.0/messages/send"

const apiUrlMailchimpSendEmailFromTemplate string = "https://mandrillapp.com/api/1.0/messages/send-template"

type IMailClient interface {
	// SendEmailVerification(ctx context.Context, tokenstring string, recipient string) error
	SendForgotPasswordOTP(ctx context.Context, otp string, recipient string) error
	// SendOnboardingGreeting(ctx context.Context, recipient string) error
}
