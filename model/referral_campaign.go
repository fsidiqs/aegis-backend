package model

import (
	"time"

	"github.com/google/uuid"
)

type ReferralCampaign struct {
	ID               uuid.UUID  `gorm:"primary_key; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	Name             string     `json:"name"`
	StartDate        *time.Time `json:"start_date"`
	EndDate          *time.Time `json:"end_date"`
	ReferrerRewardID uuid.UUID  `json:"referrer_reward_id"`
	RefereeRewardID  uuid.UUID  `json:"referee_reward_id"`
	DefaultColumns
}
