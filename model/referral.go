package model

import (
	"github.com/google/uuid"
)

type Referral struct {
	ID                 uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	ReferralCampaignID uuid.UUID `json:"referral_campaign_id"`
	ReferrerID         uuid.UUID `json:"referrer_id"`
	ReferralCode       string    `json:"referral_code"`
	DefaultColumns
}

// REQUEST DATA

type ReferralQuery struct {
	ReferralCode string `form:"referral_code" binding:"required"`
}

type ReferralReqJson struct {
	ReferralCode string `json:"referral_code"`
}
