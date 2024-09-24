package model

import (
	"github.com/fsidiqs/aegis-backend/model/apperror"
	"github.com/google/uuid"
)

type TPaymentStatus string

const (
	TStatusWaitingPayment TPaymentStatus = "WAITING_PAYMENT"
	TStatusPaymentSuccess TPaymentStatus = "SUCCESS"
	TStatusPaymentFailed  TPaymentStatus = "FAILED"
	TStatusPaymentRefund  TPaymentStatus = "REFUND"
)

type Payment struct {
	ID                 uuid.UUID       `gorm:"primary_key; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	UserID             uuid.UUID       `json:"user_id"`
	SubscriptionPlanID uuid.UUID       `json:"subscription_plan_id"`
	TxRefNumber        string          `json:"tx_ref_number"`
	Amount             int             `json:"amount"`
	ExternalRefNumber  string          `json:"external_ref_number"`
	PaymentMethod      string          `json:"payment_method"`
	PaymentChannel     string          `json:"payment_channel"`
	Status             *TPaymentStatus `json:"status"`

	DefaultColumns
}

type PaymentUpdate struct {
	PaymentMethod  string
	PaymentChannel string
	Status         *TPaymentStatus

	DefaultUpdateColumns
}

func (PaymentUpdate) TableName() string {
	return "payments"
}

//------ Request

type QueryReqPaymentStatus struct {
	Status string `form:"status"`
}

func (data QueryReqPaymentStatus) ToPaymentStatus() (*TPaymentStatus, error) {
	switch data.Status {
	case string(TStatusWaitingPayment):
		status := TStatusWaitingPayment
		return &status, nil
	case string(TStatusPaymentSuccess):
		status := TStatusPaymentSuccess
		return &status, nil
	case string(TStatusPaymentRefund):
		status := TStatusWaitingPayment
		return &status, nil
	case string(TStatusPaymentFailed):
		status := TStatusWaitingPayment
		return &status, nil

	}
	return nil, &apperror.Error{Type: apperror.BadRequest, Message: apperror.BadRequestMessage}
}
