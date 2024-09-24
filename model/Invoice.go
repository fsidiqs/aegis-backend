package model

import "time"

type InvoiceCreate struct {
	PayerEmail      string
	TXRefNum        string
	Amount          float64
	ShouldSendEmail bool
	Description     string
}

type InvoiceCreateFromEmailReq struct {
	PayerEmail     string `json:"payer_email" binding:"required,email"`
	SubscriptionID string `json:"subscription_id" binding:"required"`
}

type ExtInvoiceResp struct {
	ID                     string         `json:"id"`
	ExternalID             string         `json:"external_id"`
	UserID                 string         `json:"user_id"`
	PaymentMethod          string         `json:"payment_method"`
	Status                 TInvoiceStatus `json:"status"`
	PaidAt                 *time.Time     `json:"paid_at"`
	PayerEmail             string         `json:"payer_email"`
	Description            string         `json:"description"`
	AdjustedReceivedAmount float64        `json:"adjusted_received_amount"`
	Amount                 float64        `json:"amount"`
	PaidAmount             float64        `json:"paid_amount"`
	Updated                *time.Time     `json:"updated"`
	Created                *time.Time     `json:"created"`
	Currency               string         `json:"currency"`
	PaymentChannel         string         `json:"payment_channel"`
	PaymentDestination     string         `json:"payment_destination"`
	InvoiceURL             string         `json:"invoice_url"`
	ExpiryDate             string         `json:"expiry_date"`
}

type InvoiceResp struct {
	ID         string         `json:"id"`
	PayerEmail string         `json:"payer_email"`
	Status     TInvoiceStatus `json:"status"`
	InvoiceURL string         `json:"invoice_url,omitempty"`
}

type PaymentResp struct {
	Invoice InvoiceResp `json:"invoice"`
	Payment Payment     `json:"payment"`
}

type TInvoiceStatus string

const (
	TInvoicePaid    TInvoiceStatus = "PAID"
	TInvoiceExpired TInvoiceStatus = "EXPIRED"
	TInvoicePending TInvoiceStatus = "PENDING"
)
