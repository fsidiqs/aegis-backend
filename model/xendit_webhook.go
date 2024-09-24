package model

import (
	"time"

	"github.com/fsidiqs/aegis-backend/model/apperror"
)

type XenditWebHookReq struct {
	AdjustedReceivedAmount float64    `json:"adjusted_received_amount"`
	Amount                 float64    `json:"amount"`
	BankCode               string     `json:"bank_code"`
	Created                *time.Time `json:"created"`
	Currency               string     `json:"currency"`
	Description            string     `json:"description"`
	ExternalID             string     `json:"external_id"`
	FeesPaidAmount         float64    `json:"fees_paid_amount"`
	ID                     string     `json:"id"`
	IsHigh                 bool       `json:"is_high"`
	MerchantName           string     `json:"merchant_name"`
	PaidAmount             float64    `json:"paid_amount"`
	PaidAt                 *time.Time `json:"paid_at"`
	PayerEmail             string     `json:"payer_email"`
	PaymentChannel         string     `json:"payment_channel"`
	PaymentDestination     string     `json:"payment_destination"`
	PaymentMethod          string     `json:"payment_method"`
	Status                 string     `json:"status"`
	Updated                *time.Time `json:"updated"`
	UserID                 string     `json:"user_id"`
}

func (x *XenditWebHookReq) CreatePaymentUpdateInstance(updatedBy string) (*PaymentUpdate, error) {
	var status TPaymentStatus

	switch TInvoiceStatus(x.Status) {
	case TInvoicePaid:
		status = TStatusPaymentSuccess
	case TInvoiceExpired:
		status = TStatusPaymentFailed
	default:
		return nil, &apperror.Error{
			Type:    apperror.Internal,
			Message: apperror.MsgUnhandledPaymentStatus,
		}
	}

	return &PaymentUpdate{
		PaymentMethod:        x.PaymentMethod,
		PaymentChannel:       x.PaymentChannel,
		Status:               &status,
		DefaultUpdateColumns: NewDefaultUpdateCol(updatedBy),
	}, nil
}
