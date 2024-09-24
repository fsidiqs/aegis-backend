package model

import "github.com/google/uuid"

type MCoach struct {
	ID          uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`

	ProfilePicture       string `json:"-"`
	ProfilePictureObject string `gorm:"-" json:"profile_picture"`

	DefaultColumns
}
