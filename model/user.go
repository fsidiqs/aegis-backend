package model

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"

	helper "github.com/fsidiqs/aegis-backend/helper"
	"github.com/fsidiqs/aegis-backend/security"
)

const (
	TSUPERADMIN TRole = "superadmin"
	TUSER       TRole = "user"
)

const (
	TLastLoginManual TLastLoginMtd = "MANUAL"
)

const (
	TMale   TGender = "MALE"
	TFemale TGender = "FEMALE"
	TOther  TGender = "OTHER"
	TEmpty  TGender = ""
)

const (
	validationErrAllEmpty = "must contain atleast 1 field to update"
	hashingPasswordErr    = "error hashing password"
)

type (
	TLastLoginMtd string
	TRole         string
	TGender       string
	User          struct {
		ID              uuid.UUID       `gorm:"primary_key; unique; type:uuid; column:id; not null; default:uuid_generate_v4()" json:"id"`
		Role            TRole           `json:"role"`
		Name            string          `json:"name"`
		Email           string          `json:"email"`
		Password        string          `json:"-"`
		EmailVerifiedAt *time.Time      `json:"-"`
		LastLoginMethod string          `json:"last_login_method"`
		Organizations   []*Organization `gorm:"many2many:users_organizations;"`

		DefaultColumns
	}
)

func (u *User) AccTypeFree() {
	u.Role = TUSER
}

// ReinputForPublic updates the User Instance and hide the password value
func (u *User) ReinputForPublic(uInput *User) {
	u = &User{
		ID:       uInput.ID,
		Role:     uInput.Role,
		Name:     uInput.Name,
		Email:    uInput.Email,
		Password: "-",
		DefaultColumns: DefaultColumns{
			CreatedAt: uInput.CreatedAt,
			UpdatedAt: uInput.UpdatedAt,
		},
	}
}

func (u *User) ToUserUpdateModel() (*UserUpdate, error) {
	ret := &UserUpdate{
		User: User{
			Role: u.Role,
			Name: u.Name,
		},
	}
	return ret, nil
}

func NewGender(gender string) (TGender, bool) {
	upper := strings.ReplaceAll(strings.ToUpper(strings.Trim(gender, "")), " ", "_")
	return MatchGender(upper)
}

func MatchGender(val string) (TGender, bool) {
	switch val {
	case "MALE":
		return TMale, true
	case "FEMALE":
		return TFemale, true
	case "OTHER":
		return TOther, true
	case "":
		return TEmpty, true
	}

	return "", false
}

func NewFreeUser(name, email, role, password string) (*User, error) {
	var err error

	hPW, err := security.HashPassword(password)
	if err != nil {
		log.Printf("error instantiating a user: %v\n", err)
		return nil, errors.New(hashingPasswordErr)
	}

	var finalRole TRole

	if role == "superadmin" {
		finalRole = TSUPERADMIN
	} else {
		finalRole = TUSER
	}

	user := &User{
		Name:     name,
		Email:    email,
		Password: string(hPW),
		Role:     finalRole,
		DefaultColumns: DefaultColumns{
			CreatedBy:  helper.TraceCurrentFunc(),
			RecordFlag: RecActive,
		},
	}

	return user, nil
}

type UserUpdate struct {
	User
}

// required by GORMFramework
func (UserUpdate) TableName() string {
	return "users"
}

type UserPasswordUpdate struct {
	Password string `json:"password"`
	DefaultColumns
}

// required by GORMFramework
func (UserPasswordUpdate) TableName() string {
	return "users"
}

type UserRecordFlagUpdate struct {
	RecordFlag TRecordFlag `json:"record_flag"`
}

// required by GORMFramework
func (UserRecordFlagUpdate) TableName() string {
	return "users"
}

// REQUEST DATA

type UserRegisterReq struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required"`
	Password string `json:"password" binding:"required,gte=6,lte=30"`

	ReferralReqJson
}

const (
	validationFailedErr = "validation failed"
)

func (u UserRegisterReq) Validate() error {
	if len(u.Name) == 0 || len(u.Password) == 0 || len(u.Email) == 0 {
		return errors.New(validationFailedErr)
	}
	return nil
}

func (u UserRegisterReq) CreateModel() (*User, error) {
	err := u.Validate()
	if err != nil {
		return nil, err
	}

	userp, err := NewFreeUser(u.Name, u.Email, u.Role, u.Password)
	if err != nil {
		return nil, err
	}

	return userp, nil
}

type UserLoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6,lte=30"`
}

type UserSocialLoginReq struct {
	IDToken  string `json:"id_token" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
	Gender   string `json:"gender"`

	ReferralReqJson
}

type UserUpdateReq struct {
	User
}

func (u UserUpdateReq) Validate() error {
	return nil
}

func (u UserUpdateReq) ToUserUpdateModel(updatedBy string) (*UserUpdate, error) {
	if err := u.Validate(); err != nil {
		return nil, err
	}
	return &UserUpdate{
		User: User{
			Name: u.Name,
			Role: u.Role,
		},
	}, nil
}

type ReqUserResendVerification struct {
	Email string `json:"email" binding:"required,email"`
}

func (u ReqUserResendVerification) ToModel() (*User, error) {
	return &User{
		Email: u.Email,
	}, nil
}

type ReqUserResetPassword struct {
	Email string `json:"email" binding:"required,email"`
}

func (u ReqUserResetPassword) ToModel() (*User, error) {
	return &User{
		Email: u.Email,
	}, nil
}

// Used for endpoint /user/forgot-password/update-using-otp
type ReqUpdateForgottenPasswordByOTP struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6,lte=30"`
	OTP      string `json:"otp" binding:"required"`
}

func (u ReqUpdateForgottenPasswordByOTP) ToUserPasswordUpdateModel(updatedBy string) (*UserPasswordUpdate, error) {
	hPW, err := security.HashPassword(u.Password)
	if err != nil {
		log.Printf("error instantiating a user: %v\n", err)
		return nil, errors.New(hashingPasswordErr)
	}
	now := time.Now()
	return &UserPasswordUpdate{
		Password: string(hPW),
		DefaultColumns: DefaultColumns{
			UpdatedAt: &now,
			UpdatedBy: updatedBy,
		},
	}, nil
}

// ReqSubmitOTPForgotPassword used for /user/forgot-password/otp
type ReqSubmitOTPForgotPassword struct {
	OTP   string `json:"otp" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// Used for endpoint /user/verify-email
type ReqVerifyEmailByOTP struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required"`
}

// REQUEST PAYLOAD FOR UPGRADE A USER SUBCSRIPTION BY EMAIL

type ReqUpgradeUserSub struct {
	Email string `json:"email" binding:"required,email"`
}
