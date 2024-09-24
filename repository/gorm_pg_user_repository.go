package repository

import (
	"context"
	"fmt"

	"github.com/fsidiqs/aegis-backend/db"
	"github.com/fsidiqs/aegis-backend/helper"
	"github.com/fsidiqs/aegis-backend/helper/queryhelper"
	"github.com/fsidiqs/aegis-backend/model"
	"github.com/fsidiqs/aegis-backend/model/apperror"
	"github.com/google/uuid"
	"github.com/kr/pretty"

	"gorm.io/gorm"
)

type psqlUserRepository struct {
	DB *gorm.DB
}

// NewUserRepository is factory pattern to create UserRepository with gorm as parameter
func NewUserRepository(db *gorm.DB) IUserRepository {
	return &psqlUserRepository{
		DB: db,
	}
}

func (r *psqlUserRepository) Create(ctx context.Context, u model.User) (*model.User, error) {
	var dbContext *gorm.DB
	// get db from context, if db doesnt exists then use the class db
	if userDBCtx := ctx.Value(db.DBTrxUserKey); userDBCtx != nil {
		dbContext = userDBCtx.(*gorm.DB)
	} else {
		dbContext = r.DB
	}

	err := dbContext.WithContext(ctx).Model(&model.User{}).Create(&u).Error
	if err != nil {
		errMsg := fmt.Sprintf("%s:err:%v", helper.TraceCurrentFuncArgs(u), err)
		return nil, apperror.NewRepoErrorMsg(err, errMsg)
	}
	return &u, nil
}

func (r *psqlUserRepository) FindByID(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	var u model.User
	var dbContext *gorm.DB
	// get db from context, if db doesnt exists then use the class db
	if userDBCtx := ctx.Value(db.DBTrxUserKey); userDBCtx != nil {
		dbContext = userDBCtx.(*gorm.DB)
	} else {
		dbContext = r.DB
	}

	err := dbContext.Debug().WithContext(ctx).Where(queryhelper.ActiveFlag).
		Model(&model.User{}).First(&u, uid).Error
	if err != nil {
		errMsg := fmt.Sprintf("%s:err:%v", helper.TraceCurrentFuncArgs(
			map[string]string{
				"user_id": uid.String(),
			},
		), err)
		return nil, apperror.NewRepoErrorMsg(err, errMsg)
	}
	pretty.Println("testing", u)
	pretty.Println()
	return &u, nil
}

func (r *psqlUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var dbContext *gorm.DB

	if dbCtx := ctx.Value(db.DBTrxUserKey); dbCtx != nil {
		dbContext = dbCtx.(*gorm.DB)
	} else {
		dbContext = r.DB
	}
	var u *model.User
	err := dbContext.WithContext(ctx).Where(queryhelper.ActiveFlag).
		Model(&model.User{}).Where("email = ?", email).First(&u).Error
	if err != nil {
		errMsg := fmt.Sprintf("%s:err:%v", helper.TraceCurrentFuncArgs(
			map[string]string{
				"email": email,
			},
		), err)
		return nil, apperror.NewRepoErrorMsg(err, errMsg)
	}

	return u, nil
}

// func (r *psqlUserRepository) FindByEmailAndPhone(ctx context.Context, email string, phone string) (*model.User, error) {
// 	var dbContext *gorm.DB

// 	if dbCtx := ctx.Value(db.DBTrxUserKey); dbCtx != nil {
// 		dbContext = dbCtx.(*gorm.DB)
// 	} else {
// 		dbContext = r.DB
// 	}
// 	var u *model.User
// 	err := dbContext.WithContext(ctx).
// 		Model(&model.User{}).Where("email = ? AND phone_number = ?", email, phone).First(&u).Error

// 	if errors.Is(err, gorm.ErrRecordNotFound) {
// 		return nil, nil
// 	}

// 	if err != nil {
// 		errMsg := fmt.Sprintf("%s:err:%v", helper.TraceCurrentFuncArgs(
// 			map[string]string{
// 				"email": email,
// 				"phone": phone,
// 			},
// 		), err)
// 		return nil, apperror.NewRepoErrorMsg(err, errMsg)
// 	}
// 	return u, nil
// }

func (r *psqlUserRepository) Update(ctx context.Context, uid uuid.UUID, u *model.UserUpdate) error {
	var dbContext *gorm.DB

	if dbCtx := ctx.Value(db.DBTrxUserKey); dbCtx != nil {
		dbContext = dbCtx.(*gorm.DB)
	} else {
		dbContext = r.DB
	}

	rows := dbContext.WithContext(ctx).
		Where(queryhelper.ActiveFlag).
		// Model(&model.User{}).
		Where("id = ?", uid).
		Updates(u).
		RowsAffected

	if rows == 0 {
		errMsg := fmt.Sprintf("%s", helper.TraceCurrentFuncArgs(
			map[string]string{
				"user_id": uid.String(),
			},
			*u,
		))
		return apperror.NewRepoErrorMsg(nil, errMsg)
	}
	return nil
}

func (r *psqlUserRepository) UpdatePassword(ctx context.Context, uid uuid.UUID, u model.UserPasswordUpdate) error {
	err := r.DB.WithContext(ctx).
		Where(queryhelper.ActiveFlag).
		Where("id = ?", uid).
		Updates(u).
		Error
	if err != nil {
		errMsg := fmt.Sprintf("%s", helper.TraceCurrentFuncArgs(
			map[string]string{
				"user_id": uid.String(),
			},
			u,
		))
		return apperror.NewRepoErrorMsg(err, errMsg)
	}
	return nil
}

func (r *psqlUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	var dbContext *gorm.DB
	// get db from context, if db doesnt exists then use the class db
	if userDBCtx := ctx.Value(db.DBTrxUserKey); userDBCtx != nil {
		dbContext = userDBCtx.(*gorm.DB)
	} else {
		dbContext = r.DB
	}

	var ex bool
	// check if user_id AND program_id is already exists
	err := dbContext.WithContext(ctx).Raw(queryhelper.EmailExists, email).
		Find(&ex).
		Error
	if err != nil {
		errMsg := fmt.Sprintf("%s", helper.TraceCurrentFuncArgs(
			map[string]string{
				"email": email,
			},
		))
		return false, apperror.NewRepoErrorMsg(err, errMsg)
	}
	return ex, nil
}

func (r *psqlUserRepository) StoreSession(ctx context.Context, s *model.UserSession) error {
	var dbContext *gorm.DB

	if dbCtx := ctx.Value(db.DBTrxUserKey); dbCtx != nil {
		dbContext = dbCtx.(*gorm.DB)
	} else {
		dbContext = r.DB
	}

	err := dbContext.WithContext(ctx).
		Model(&model.UserSession{}).Create(&s).Error
	if err != nil {
		errMsg := fmt.Sprintf("%s", helper.TraceCurrentFuncArgs(*s))
		return apperror.NewRepoErrorMsg(err, errMsg)
	}

	return nil
}

// func (r *psqlUserRepository) FindUserSessionDBByUserID(ctx context.Context, uid uuid.UUID) (*model.UserSession, error) {
// 	var dbContext *gorm.DB

// 	if dbCtx := ctx.Value(db.DBTrxUserKey); dbCtx != nil {
// 		dbContext = dbCtx.(*gorm.DB)
// 	} else {
// 		dbContext = r.DB
// 	}
// 	var uSess *model.UserSession
// 	err := dbContext.WithContext(ctx).
// 		Model(&model.UserSession{}).
// 		Where("user_id = ?", uid).
// 		Where(queryhelper.ActiveFlag).
// 		First(&uSess).Error
// 	if err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			return nil, nil
// 		}
// 		errMsg := fmt.Sprintf("%s", helper.TraceCurrentFuncArgs(
// 			map[string]string{
// 				"user_id": uid.String(),
// 			},
// 		))
// 		return nil, apperror.NewRepoErrorMsg(err, errMsg)
// 	}

// 	return uSess, nil
// }

func (r *psqlUserRepository) RemoveSession(ctx context.Context, userID uuid.UUID, prevrefreshTokenID uuid.UUID) error {
	var dbContext *gorm.DB

	if dbCtx := ctx.Value(db.DBTrxUserKey); dbCtx != nil {
		dbContext = dbCtx.(*gorm.DB)
	} else {
		dbContext = r.DB
	}
	err := dbContext.WithContext(ctx).
		Where(queryhelper.ActiveFlag).
		Where("user_id = ? AND refresh_token_id = ?", userID, prevrefreshTokenID).
		Delete(&model.UserSession{}).Error
	if err != nil {
		errMsg := fmt.Sprintf("%s", helper.TraceCurrentFuncArgs(
			map[string]string{
				"user_id":            userID.String(),
				"prev_refresh_token": prevrefreshTokenID.String(),
			},
		))
		return apperror.NewRepoErrorMsg(err, errMsg)
	}

	return nil
}

func (r *psqlUserRepository) RemoveUserSessions(ctx context.Context, userID uuid.UUID) error {
	var dbContext *gorm.DB

	if dbCtx := ctx.Value(db.DBTrxUserKey); dbCtx != nil {
		dbContext = dbCtx.(*gorm.DB)
	} else {
		dbContext = r.DB
	}
	err := dbContext.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&model.UserSession{}).Error
	if err != nil {
		errMsg := fmt.Sprintf("%s", helper.TraceCurrentFuncArgs(
			map[string]string{
				"user_id": userID.String(),
			},
		))
		return apperror.NewRepoErrorMsg(err, errMsg)
	}

	return nil
}

func (r *psqlUserRepository) UpdateUserSession(ctx context.Context, userID uuid.UUID, uSessUpdate model.UserSessionUpdate) error {
	var dbContext *gorm.DB

	if dbCtx := ctx.Value(db.DBTrxUserKey); dbCtx != nil {
		dbContext = dbCtx.(*gorm.DB)
	} else {
		dbContext = r.DB
	}
	err := dbContext.WithContext(ctx).
		Where(queryhelper.ActiveFlag).
		Where("user_id = ?", userID).
		Updates(&uSessUpdate).
		Error
	if err != nil {
		errMsg := fmt.Sprintf("%s", helper.TraceCurrentFuncArgs(
			map[string]string{
				"user_id": userID.String(),
			},
			uSessUpdate,
		))
		return apperror.NewRepoErrorMsg(err, errMsg)
	}
	return err
}

func (r *psqlUserRepository) RefreshTokenExists(ctx context.Context, uid uuid.UUID, refreshTokenID uuid.UUID) (bool, error) {
	var dbContext *gorm.DB

	if dbCtx := ctx.Value(db.DBTrxUserKey); dbCtx != nil {
		dbContext = dbCtx.(*gorm.DB)
	} else {
		dbContext = r.DB
	}
	err := dbContext.WithContext(ctx).
		Where("refresh_token_id = ? AND user_id = ? AND status = ?", refreshTokenID.String(), uid, "ACTIVE").
		Where(queryhelper.ActiveFlag).
		First(&model.UserSession{}).Error
	if err != nil {
		errMsg := fmt.Sprintf("%s", helper.TraceCurrentFuncArgs(
			map[string]string{
				"user_id":          uid.String(),
				"refresh_token_id": refreshTokenID.String(),
			},
		))
		return false, apperror.NewRepoErrorMsg(err, errMsg)
	}
	return true, nil
}

func (r *psqlUserRepository) FindOTP(ctx context.Context, userID uuid.UUID, otp string, otpType model.TOTP) (*model.UserOTP, error) {
	var dbContext *gorm.DB

	if dbCtx := ctx.Value(db.DBTrxUserKey); dbCtx != nil {
		dbContext = dbCtx.(*gorm.DB)
	} else {
		dbContext = r.DB
	}
	var uOTP *model.UserOTP
	err := dbContext.WithContext(ctx).Model(&model.UserOTP{}).
		Where("user_id = ? AND otp = ? AND type = ?", userID, otp, otpType).First(&uOTP).Error
	if err != nil {
		errMsg := fmt.Sprintf("%s", helper.TraceCurrentFuncArgs(
			map[string]string{
				"user_id": userID.String(),
				"otp":     otp,
			},
			otpType,
		))
		return nil, apperror.NewRepoErrorMsg(err, errMsg)
	}

	if uOTP == nil {
		errMsg := fmt.Sprintf("%s", helper.TraceCurrentFuncArgs(
			map[string]string{
				"user_id": userID.String(),
				"otp":     otp,
			},
			otpType,
		))
		return nil, apperror.NewRepoErrorMsg(err, errMsg)
	}
	return uOTP, nil
}

func (r *psqlUserRepository) UpsertOTP(ctx context.Context, userID string, uOTP model.UserOTP) error {
	// var err error

	// var dbContext *gorm.DB
	// get db from context, if db doesnt exists then use the class db
	// if userDBCtx := ctx.Value(db.DBTrxUserKey); userDBCtx != nil {
	// 	dbContext = userDBCtx.(*gorm.DB)
	// } else {
	// 	dbContext = r.DB
	// }

	// now := time.Now()
	// uOTPUpdate := model.UserOTPUpdate{
	// 	OTP:       uOTP.OTP,
	// 	Type:      uOTP.Type,
	// 	Status:    uOTP.Status,
	// 	ExpiredAt: uOTP.ExpiredAt,
	// 	DefaultColumns: model.DefaultColumns{
	// 		UpdatedAt: &now,
	// 		UpdatedBy: helper.TraceCurrentFunc(),
	// 	},
	// }
	// if dbContext.Model(&model.UserOTPUpdate{}).Where("user_id = ? AND type = ?", userID, uOTP.Type).Updates(&uOTPUpdate).RowsAffected == 0 {
	// 	err = dbContext.Create(&uOTP).Error
	// 	if err != nil {
	// 		errMsg := fmt.Sprintf("%s", helper.TraceCurrentFuncArgs(
	// 			map[string]string{
	// 				"user_id": userID,
	// 			},
	// 			uOTP,
	// 		))
	// 		return apperror.NewRepoErrorMsg(err, errMsg)
	// 	}
	// }
	return nil
}

func (r *psqlUserRepository) UpdateOTP(ctx context.Context, uOTPID uuid.UUID, userOTP model.UserOTPUpdate) error {
	var dbContext *gorm.DB

	if dbCtx := ctx.Value(db.DBTrxUserKey); dbCtx != nil {
		dbContext = dbCtx.(*gorm.DB)
	} else {
		dbContext = r.DB
	}
	err := dbContext.WithContext(ctx).
		Where(queryhelper.ActiveFlag).
		// Model(&model.User{}).
		Where("id = ?", uOTPID).
		Updates(&userOTP).
		Error
	if err != nil {
		errMsg := fmt.Sprintf("%s", helper.TraceCurrentFuncArgs(
			map[string]string{
				"user_otp_id": uOTPID.String(),
			},
			userOTP,
		))
		return apperror.NewRepoErrorMsg(err, errMsg)
	}
	return nil
}

// func (r *psqlUserRepository) GetSocialLoginBySocialUserID(ctx context.Context, socUserID string) (*model.SocialLogin, error) {
// 	var dbContext *gorm.DB
// 	// get db from context, if db doesnt exists then use the class db
// 	if userDBCtx := ctx.Value(db.DBTrxUserKey); userDBCtx != nil {
// 		dbContext = userDBCtx.(*gorm.DB)
// 	} else {
// 		dbContext = r.DB
// 	}
// 	var socLogin model.SocialLogin

// 	err := dbContext.WithContext(ctx).
// 		Model(&model.SocialLogin{}).
// 		Where(queryhelper.ActiveFlag).
// 		Where("social_user_id = ?", socUserID).
// 		First(&socLogin).Error
// 	if err != nil {
// 		errMsg := fmt.Sprintf("%s", helper.TraceCurrentFuncArgs(
// 			map[string]string{
// 				"social_user_id": socUserID,
// 			}))
// 		return nil, apperror.NewRepoErrorMsg(err, errMsg)
// 	}

// 	return &socLogin, nil
// }

// func (r *psqlUserRepository) SocialLoginExists(ctx context.Context, socialUserID string, provider string) (bool, error) {
// 	var dbContext *gorm.DB

// 	if dbCtx := ctx.Value(db.DBTrxUserKey); dbCtx != nil {
// 		dbContext = dbCtx.(*gorm.DB)
// 	} else {
// 		dbContext = r.DB
// 	}
// 	var ex bool
// 	// check if user_id AND program_id is already exists
// 	err := dbContext.WithContext(ctx).Raw(queryhelper.SocialLoginExists, socialUserID, provider).
// 		Find(&ex).
// 		Error
// 	if err != nil {
// 		errMsg := fmt.Sprintf("%s", helper.TraceCurrentFuncArgs(
// 			map[string]string{
// 				"social_user_id": socialUserID,
// 				"provider":       provider,
// 			},
// 		))
// 		return false, apperror.NewRepoErrorMsg(err, errMsg)
// 	}
// 	return ex, nil
// }

// func (r *psqlUserRepository) SocialLoginStore(ctx context.Context, u model.SocialLogin) (*model.SocialLogin, error) {
// 	var dbContext *gorm.DB

// 	if dbCtx := ctx.Value(db.DBTrxUserKey); dbCtx != nil {
// 		dbContext = dbCtx.(*gorm.DB)
// 	} else {
// 		dbContext = r.DB
// 	}
// 	var err error

// 	err = dbContext.Create(&u).Error
// 	if err != nil {
// 		errMsg := fmt.Sprintf("%s", helper.TraceCurrentFuncArgs(u))
// 		return nil, apperror.NewRepoErrorMsg(err, errMsg)
// 	}
// 	return &u, nil
// }

// func (r *psqlUserRepository) SocialLoginUpdate(ctx context.Context, socLoginID uuid.UUID, socialLoginUp model.SocialLoginUpdate) error {
// 	var dbContext *gorm.DB

// 	if dbCtx := ctx.Value(db.DBTrxUserKey); dbCtx != nil {
// 		dbContext = dbCtx.(*gorm.DB)
// 	} else {
// 		dbContext = r.DB
// 	}
// 	err := dbContext.WithContext(ctx).
// 		Where(queryhelper.ActiveFlag).
// 		Where("id = ?", socLoginID).
// 		Updates(&socialLoginUp).
// 		Error
// 	if err != nil {
// 		errMsg := fmt.Sprintf("%s", helper.TraceCurrentFuncArgs(
// 			map[string]string{
// 				"social_login_id": socLoginID.String(),
// 			},
// 			socialLoginUp,
// 		))
// 		return apperror.NewRepoErrorMsg(err, errMsg)
// 	}

// 	return nil
// }

// func (r *psqlUserRepository) UserUpdateByEmail(ctx context.Context, email string, u model.UserUpdate) error {
// 	var dbContext *gorm.DB

// 	if dbCtx := ctx.Value(db.DBTrxUserKey); dbCtx != nil {
// 		dbContext = dbCtx.(*gorm.DB)
// 	} else {
// 		dbContext = r.DB
// 	}
// 	err := dbContext.WithContext(ctx).
// 		Where(queryhelper.ActiveFlag).
// 		// Model(&model.User{}).
// 		Where("email = ?", email).
// 		Updates(&u).
// 		Error
// 	if err != nil {
// 		errMsg := fmt.Sprintf("%s", helper.TraceCurrentFuncArgs(
// 			map[string]string{
// 				"email": email,
// 			},
// 			u,
// 		))
// 		return apperror.NewRepoErrorMsg(err, errMsg)
// 	}
// 	return nil
// }

// func (r *psqlUserRepository) IsSubscribing(ctx context.Context, userID uuid.UUID) (bool, error) {
// 	var dbContext *gorm.DB

// 	if dbCtx := ctx.Value(db.DBTrxUserKey); dbCtx != nil {
// 		dbContext = dbCtx.(*gorm.DB)
// 	} else {
// 		dbContext = r.DB
// 	}
// 	var isSubsribing bool
// 	now := time.Now()
// 	err := dbContext.WithContext(ctx).
// 		Raw(queryhelper.IsUserSubscribing, userID, now).
// 		Find(&isSubsribing).Error
// 	if err != nil {
// 		errMsg := fmt.Sprintf("%s", helper.TraceCurrentFuncArgs(
// 			map[string]string{
// 				"user_id": userID.String(),
// 			},
// 		))
// 		return false, apperror.NewRepoErrorMsg(err, errMsg)
// 	}
// 	return isSubsribing, nil
// }
