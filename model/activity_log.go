package model

import (
	"time"

	"github.com/google/uuid"
)

type UserActivityLog struct {
	ID                uuid.UUID  `gorm:"primary_key; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	UserId            uuid.UUID  `json:"user_id"`
	RelatedActivityID *uuid.UUID `json:"related_activity_id"`
	ActivityDate      time.Time  `json:"activity_date"`
	ActivityType      string     `json:"type"`

	DefaultColumns
}
