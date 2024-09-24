package model

import (
	"github.com/google/uuid"
)

type CampaignReward struct {
	ID         uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	Name       string    `json:"name"`
	RewardType string    `json:"reward_type"`
	Amount     int       `json:"amount"`

	DefaultColumns
}
