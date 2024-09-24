package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fsidiqs/aegis-backend/mail"
)

type mailer struct {
	client mail.IMailClient
}

func initMailer() (*mailer, error) {
	var err error
	var appDomain, apiKey, fromName, fromEmail, port, apiUrl string
	var mailerClient mail.IMailClient

	apiKey = os.Getenv("MAILCHIMP_API_KEY")
	fromName = os.Getenv("MAIL_FROM_NAME")
	fromEmail = os.Getenv("MAIL_FROM_EMAIL")
	appDomain = os.Getenv("APP_DOMAIN")
	port = os.Getenv("PORT")
	apiUrl = os.Getenv("API_URL")

	log.Printf("Initializing mailer\n")

	mailerClient, err = mail.NewMailChimpMailable(
		apiKey,
		fromEmail,
		fromName,
		fmt.Sprintf("%s:%s%s", appDomain, port, apiUrl),
	)
	if err != nil {
		return nil, fmt.Errorf("error connecting creating mail instance: %v", err)
	}
	log.Printf("mailer has created")

	return &mailer{
		client: mailerClient,
	}, nil
}
