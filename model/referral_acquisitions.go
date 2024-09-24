package model

import (
	"github.com/google/uuid"
)

type ReferralAcquisition struct {
	ID         uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	ReferralID uuid.UUID `json:"-"`
	RefereeID  uuid.UUID `json:"-"`
	DeviceID   string    `json:"-"`
	DefaultColumns
}

type ReferralAcquisitionUpdate struct {
	RefereeID uuid.UUID
	DefaultUpdateColumns
}
