package model

import "github.com/google/uuid"

type MSafeRoom struct {
	ID              uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	VendorRoomID    uuid.UUID `json:"vendor_room_id"`
	Name            string    `json:"name"`
	ModeratorUserID uuid.UUID `json:"moderator_user_id"`

	DefaultColumns
}
