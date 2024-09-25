package model

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID          uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	Name        string    `gorm:"size:255;not null"` // Name of the organization, with a max size of 255
	Description string    `gorm:"type:text"`         // Detailed description (optional)
	CreatorID   uuid.UUID `gorm:"not null"`          // Reference to the user who created the organization
	CreatedAt   time.Time `gorm:"autoCreateTime"`    // Automatically set to current time on creation
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`    // Automatically set to current time on update

	// Associations
	Creator User `gorm:"foreignKey:CreatorID;constraint:OnDelete:CASCADE"` // Establishes relationship with the User model
}
