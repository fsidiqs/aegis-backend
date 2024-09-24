package model

import (
	"time"

	"github.com/google/uuid"
)

type MSubscriptionPlan struct {
	ID            uuid.UUID  `gorm:"primary_key; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	Name          string     `json:"name"`
	Price         int        `json:"price"`
	Duration      int        `json:"duration"`
	OriginalPrice int        `json:"original_price"`
	DiscountEndAt *time.Time `json:"discount_end_at"`
	DefaultColumns
}

type MSubscriptionPlanResponse struct {
	MSubscriptionPlan
	HasDiscount bool `json:"has_discount"`
}
