package model

import "github.com/google/uuid"

const (
	TSocialLoginGoogleProvider = "google.com"
)

type SocialLogin struct {
	ID           uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	SocialToken  string    `json:"-"`
	SocialUserID string    `json:"social_user_id"`
	Provider     string    `json:"provider"`

	DefaultColumns
}

type TSocialloginProvider string

type SocialLoginReq struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	SocialUserID string `json:"social_user_id"`
	Provider     string `json:"provider"`
}

func (d *SocialLoginReq) ValidateAndFix(nickname string) {
	if len(d.Name) == 0 {
		d.Name = nickname
	}
}

type SocialLoginUpdate struct {
	SocialToken  string `json:"-"`
	SocialUserID string `json:"social_user_id"`
	Provider     string `json:"provider"`

	DefaultColumns
}

func (SocialLoginUpdate) TableName() string {
	return "social_logins"
}
