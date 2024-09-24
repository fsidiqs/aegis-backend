package model

import (
	"github.com/google/uuid"
)

type EmailVerificationToken struct {
	ID      uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	TokenID uuid.UUID
	UserID  uuid.UUID

	DefaultColumns
}
