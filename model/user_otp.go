package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	TotpForgotPassword    TOTP = "FORGOT_PASSWORD"
	TotpRegistVerifyEmail TOTP = "VERIFY_EMAIL"
)

const (
	TOTPCreated        TOTPStatus = "CREATED"
	TOTPVerified       TOTPStatus = "VERIFIED"
	TOTPStatusFinished TOTPStatus = "FINISHED"
)

type UserOTP struct {
	ID        uuid.UUID `gorm:"primary_key; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID
	OTP       string
	Type      TOTP
	Status    TOTPStatus
	ExpiredAt time.Time `json:"expired_at"`

	DefaultColumns
}

func (UserOTP) TableName() string {
	return "user_otps"
}

type TOTP string
type TOTPStatus string

type UserOTPUpdate struct {
	OTP       string
	Type      TOTP
	Status    TOTPStatus
	ExpiredAt time.Time `json:"expired_at"`

	DefaultColumns
}

// required by GORMFramework
func (UserOTPUpdate) TableName() string {
	return "user_otps"
}
