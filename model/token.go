package model

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken stores token properties that
// are accessed in multiple application layers
//
// swagger:model RefreshToken
type RefreshToken struct {
	ID uuid.UUID `json:"-"`
	// user_id
	UID       uuid.UUID     `json:"-"`
	SS        string        `json:"refresh_token"`
	ExpiresIn time.Duration `json:"-"`
}

// AuthToken stores token properties that
// are accessed in multiple application layers
//
// swagger:model AuthToken

type AuthToken struct {
	SS string `json:"auth_token"`
}

// TokenData used for returning pairs of id and refresh tokens
//
// swagger:response TokenData
type TokenData struct {
	AuthToken
	RefreshToken
	NotificationToken string `json:"notification_token,omitempty"`
}

type PublicToken struct {
	SS string `json:"public_token"`
}

type PublicTokenData struct {
	PublicToken
}

type ValidatedAuthData struct {
	User           *User  `json:"user"`
	RefreshTokenID string `json:"-"`
	Provider       string `json:"auth_provider"`
}

type ValidatedPublicToken struct {
	DeviceID          string `json:"device_id"`
	NotificationToken string `json:"notification_token"`
}

// EmailVerificationTokenData used for email verification purpose

type EmailVerificationTokenData struct {
	ID        uuid.UUID     `json:"-"`
	SS        string        `json:"email_verificationToken"`
	ExpiresIn time.Duration `json:"-"`
}

// OTP
type OTPData struct {
	OTP       string    `json:"otp"`
	ExpiresAt time.Time `json:"-"`
}

// REQUEST HANDLER -----------------------------
type TokensReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
