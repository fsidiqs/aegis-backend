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

func initMailer() (mail.IMailClient, error) {
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

	return mailerClient, nil
}

func initSendgridMailer() (mail.IMailClient, error) {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	fromName := os.Getenv("MAIL_FROM_NAME")
	fromEmail := os.Getenv("MAIL_FROM_EMAIL")
	fmt.Println("Initializing mailer\n")
	fmt.Println("apiKey: ", apiKey)
	mailerClient, err := mail.NewSendGridMailable(apiKey, fromEmail, fromName)
	if err != nil {
		return nil, err
	}
	return mailerClient, nil
}

func initMailersend() (mail.IMailClient, error) {
	apiKey := os.Getenv("MAILERSEND_API_KEY")
	fromName := os.Getenv("MAIL_FROM_NAME")
	fromEmail := os.Getenv("MAIL_FROM_EMAIL")
	fmt.Println("Initializing mailer\n")
	fmt.Println("apiKey: ", apiKey)
	mailerClient, err := mail.NewMailersendMailable(apiKey, fromEmail, fromName)
	if err != nil {
		return nil, err
	}
	return mailerClient, nil
}
