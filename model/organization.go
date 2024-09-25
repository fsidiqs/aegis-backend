package model

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID          uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`    // Name of the organization, with a max size of 255
	Description string    `gorm:"type:text" json:"description"`     // Detailed description (optional)
	CreatorID   uuid.UUID `gorm:"not null" json:"creator_id"`       // Reference to the user who created the organization
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"` // Automatically set to current time on creation
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"` // Automatically set to current time on update

	// Associations
	Creator User `gorm:"foreignKey:CreatorID;constraint:OnDelete:CASCADE"` // Establishes relationship with the User model
}
