package mail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type sendgridImpl struct {
	sendGridClient      *sendgrid.Client
	fromEmail           string
	fromName            string
	apiKey              string
	onboardingFromEmail string
	onboardingFromName  string
}

func NewSendGridMailable(apiKey string, fromEmail string, fromName string) (IMailClient, error) {
	client := sendgrid.NewSendClient(apiKey)

	mail := &sendgridImpl{
		sendGridClient: client,
		fromEmail:      fromEmail,
		fromName:       fromName,
		apiKey:         apiKey,
	}
	return mail, nil
}

func (m *sendgridImpl) makeMailchimpMail(html string, plainTextContent string, subject string, recipients []map[string]interface{}) []byte {
	message := map[string]interface{}{
		"html":                     html,
		"text":                     plainTextContent,
		"subject":                  subject,
		"from_email":               m.fromEmail,
		"from_name":                m.fromName,
		"to":                       recipients,
		"important":                false,
		"track_opens":              false,
		"track_clicks":             false,
		"auto_text":                false,
		"auto_html":                false,
		"inline_css":               false,
		"url_strip_qs":             false,
		"preserve_recipients":      false,
		"view_content_link":        false,
		"merge":                    false,
		"merge_language":           "mailchimp",
		"global_merge_vars":        []string{},
		"merge_vars":               []string{},
		"tags":                     []string{},
		"google_analytics_domains": []string{},
		"recipient_metadata":       []string{},
		"attachments":              []string{},
		"images":                   []string{},
	}

	reqBody, _ := json.Marshal(map[string]interface{}{
		"key":     m.apiKey,
		"message": message,
		"async":   false,
		"ip_pool": "",
		"send_at": "",
	})
	return reqBody
}

func (m *sendgridImpl) makeMailChimpFromTemplate(recipients []map[string]interface{}) []byte {
	message := map[string]interface{}{
		// "text":         plainTextContent,
		"subject":      "Yeay! Selamat Bergabung di Mindtera!",
		"from_email":   m.onboardingFromEmail,
		"from_name":    m.onboardingFromName,
		"to":           recipients,
		"important":    false,
		"track_opens":  false,
		"track_clicks": false,
	}

	reqBody, _ := json.Marshal(map[string]interface{}{
		"key":           m.apiKey,
		"message":       message,
		"template_name": "template-app-onboarding-new-user",
		"template_content": []map[string]string{
			{"name": "", "content": ""},
		},
	})
	return reqBody
}

// func (m *sendgridImpl) SendEmailVerification(ctx context.Context, otp, recipient string) error {
// 	htmlContent := fmt.Sprintf("<p>Kode OTP verifikasi akun Anda adalah %v</p>", otp)
// 	plainTextContent := fmt.Sprintf("Kode OTP verifikasi akun Anda adalah %v", otp)
// 	subject := "[MINDTERA] Kode Verifikasi Akun"
// 	recipients := []map[string]interface{}{
// 		{
// 			"email": recipient,
// 		},
// 	}
// 	reqBody := m.makeMailchimpMail(htmlContent, plainTextContent, subject, recipients)
// 	resp, err := http.Post(apiUrlMailchimpSendEmail, applicationJsonContentType, bytes.NewBuffer(reqBody))
// 	if err != nil {
// 		log.Printf("failed to send verification email:%v \n", err)
// 		return err
// 	}
// 	respBody, _ := ioutil.ReadAll(resp.Body)
// 	var parsedResp interface{}
// 	err = json.Unmarshal(respBody, &parsedResp)
// 	if err != nil {
// 		log.Printf("failed to parse email:%v \n", err)
// 	}

// 	return nil
// }

func (m *sendgridImpl) SendOnboardingGreeting(ctx context.Context, recipient string) error {
	recipients := []map[string]interface{}{
		{
			"email": recipient,
		},
	}
	reqBody := m.makeMailChimpFromTemplate(recipients)
	resp, err := http.Post(apiUrlMailchimpSendEmailFromTemplate, applicationJsonContentType, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("failed to send greeting email:%v \n", err)
		return err
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	var parsedResp interface{}
	err = json.Unmarshal(respBody, &parsedResp)
	if err != nil {
		log.Printf("failed to parse greeting email response:%v \n", err)
	}

	return nil
}

func (m *sendgridImpl) SendForgotPasswordOTP(ctx context.Context, otp, recipient string) error {
	htmlContent := fmt.Sprintf("<p>Kode OTP reset password Anda adalah %v</p>", otp)
	plainTextContent := fmt.Sprintf("Kode OTP reset password Anda adalah %v", otp)
	subject := "[AEGIS-BACKEND-TEST] Kode Verifikasi Reset Password"
	// reqBody := m.makeMailchimpMail(htmlContent, plainTextContent, subject, recipients)

	// _, err := http.Post(apiUrlMailchimpSendEmail, applicationJsonContentType, bytes.NewBuffer(reqBody))
	// if err != nil {
	// 	log.Printf("failed to send forgot_password email:%v \n", err)
	// 	return err
	// }
	fromEmail := mail.NewEmail(m.fromName, m.fromEmail)
	recipientEmail := mail.NewEmail(recipient, recipient)
	message := mail.NewSingleEmail(
		fromEmail,
		subject,
		recipientEmail,
		plainTextContent,
		htmlContent,
	)
	response, err := m.sendGridClient.Send(message)
	if err != nil {
		log.Printf("failed to send forgot_password email:%v \n", err)
		return err
	}
	log.Printf("response: %v", response)
	return nil
}

func (m *sendgridImpl) SendAccountCreatedMail(ctx context.Context, password string, recipient string) error {
	htmlContent := fmt.Sprintf("<p>Hai %s, password Anda adalah %v</p>", recipient, password)
	plainTextContent := fmt.Sprintf("Hai %s, password Anda adalah %v", recipient, password)
	subject := "[AEGIS-BACKEND-TEST] Akun Telah Dibuat"
	// reqBody := m.makeMailchimpMail(htmlContent, plainTextContent, subject, recipients)

	// _, err := http.Post(apiUrlMailchimpSendEmail, applicationJsonContentType, bytes.NewBuffer(reqBody))
	// if err != nil {
	// 	log.Printf("failed to send forgot_password email:%v \n", err)
	// 	return err
	// }
	fromEmail := mail.NewEmail(m.fromName, m.fromEmail)
	recipientEmail := mail.NewEmail(recipient, recipient)
	message := mail.NewSingleEmail(
		fromEmail,
		subject,
		recipientEmail,
		plainTextContent,
		htmlContent,
	)
	response, err := m.sendGridClient.Send(message)
	if err != nil {
		log.Printf("failed to send forgot_password email:%v \n", err)
		return err
	}
	log.Printf("response: %v", response)
	return nil
}
