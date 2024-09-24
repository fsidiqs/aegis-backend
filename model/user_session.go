package model

import (
	"encoding/json"
	"time"

	"github.com/fsidiqs/aegis-backend/helper"
	"github.com/google/uuid"
)

const (
	TUserSessionActive TUserSessionStatus = "ACTIVE"
	TUserSessionLogout TUserSessionStatus = "LOGOUT"
)

// swagger:model UserSession
type UserSession struct {
	ID                uuid.UUID          `gorm:"primary_key; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
	AuthToken         string             `json:"auth_token"`
	RefreshTokenID    string             `json:"refresh_token_id"`
	NotificationToken string             `json:"notification_token"`
	ExpiredAt         time.Time          `json:"expired_at"`
	Status            TUserSessionStatus `json:"status"`
	UserID            uuid.UUID          `json:"user_id"`
	Role              TRole              `gorm:"-" json:"account_type"`
	UserRecordFlag    TRecordFlag        `gorm:"-" json:"user_record_flag"`

	DefaultColumns
}

type TUserSessionStatus string

type RedisUserSessionVal struct {
	Role           TRole       `gorm:"-" json:"account_type"`
	UserRecordFlag TRecordFlag `json:"user_record_flag"`
}

func (u *UserSession) MarshalJSON() ([]byte, error) {
	return json.Marshal(&RedisUserSessionVal{
		Role:           u.Role,
		UserRecordFlag: u.UserRecordFlag,
	})
}

type ExtraData struct {
	NotificationToken string
	Status            TUserSessionStatus
	Role              TRole
	UserRecordFlag    TRecordFlag
}

func NewUserSession(td *TokenData, extra ExtraData) UserSession {
	currentT := time.Now()
	exp := currentT.Add(td.ExpiresIn)
	return UserSession{
		UserID:            td.RefreshToken.UID,
		Role:              extra.Role,
		AuthToken:         td.AuthToken.SS,
		Status:            extra.Status,
		RefreshTokenID:    td.RefreshToken.ID.String(),
		NotificationToken: extra.NotificationToken,
		ExpiredAt:         exp,
		UserRecordFlag:    extra.UserRecordFlag,
		DefaultColumns: DefaultColumns{
			CreatedBy:  helper.TraceCurrentFunc(),
			RecordFlag: RecActive,
		},
	}
}

type UserSessionUpdate struct {
	AuthToken         string             `json:"auth_token"`
	RefreshTokenID    string             `json:"refresh_token_id"`
	NotificationToken string             `json:"notification_token"`
	ExpiredAt         time.Time          `json:"expired_at"`
	Status            TUserSessionStatus `json:"status"`
	Role              TRole              `gorm:"-" json:"account_type"`
	UserRecordFlag    TRecordFlag        `gorm:"-" json:"user_record_flag"`

	DefaultColumns
}

func (UserSessionUpdate) TableName() string {
	return "user_sessions"
}

func (d UserSessionUpdate) ToUserSession(uid uuid.UUID) (*UserSession, error) {
	return &UserSession{
		AuthToken:         d.AuthToken,
		RefreshTokenID:    d.RefreshTokenID,
		NotificationToken: d.NotificationToken,
		ExpiredAt:         d.ExpiredAt,
		Status:            d.Status,
		UserID:            uid,
		Role:              d.Role,
		UserRecordFlag:    d.UserRecordFlag,
		DefaultColumns: DefaultColumns{
			CreatedAt:  d.CreatedAt,
			UpdatedAt:  d.UpdatedAt,
			RecordFlag: d.RecordFlag,
			CreatedBy:  d.CreatedBy,
			UpdatedBy:  d.UpdatedBy,
		},
	}, nil
}

func NewUserSessionUpdate(td *TokenData, extra ExtraData) UserSessionUpdate {
	currentT := time.Now()
	exp := currentT.Add(td.ExpiresIn)
	return UserSessionUpdate{
		AuthToken:         td.AuthToken.SS,
		RefreshTokenID:    td.RefreshToken.ID.String(),
		NotificationToken: extra.NotificationToken,
		ExpiredAt:         exp,
		DefaultColumns: DefaultColumns{
			CreatedBy:  helper.TraceCurrentFunc(),
			RecordFlag: RecActive,
		},
	}
}
