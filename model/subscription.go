package model

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID                 uuid.UUID          `gorm:"primary_key; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	StartDate          *time.Time         `json:"start_date"`
	EndDate            *time.Time         `json:"end_date"`
	BilledAmount       int                `json:"billed_amount"`
	UserID             uuid.UUID          `json:"user_id"`
	PaymentID          uuid.UUID          `json:"payment_id"`
	SubscriptionPlanID uuid.UUID          `json:"subscription_plan_id"`
	SubscriptionPlan   *MSubscriptionPlan `gorm:"foreignKey:SubscriptionPlanID" json:"subscription_plan"`
	DefaultColumns
}
