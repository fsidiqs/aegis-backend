package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/fsidiqs/aegis-backend/helper"
	"github.com/fsidiqs/aegis-backend/mail"
	"github.com/fsidiqs/aegis-backend/model"
	"github.com/fsidiqs/aegis-backend/model/apperror"
	"github.com/fsidiqs/aegis-backend/repository"
	"github.com/fsidiqs/aegis-backend/security"
)

// UserService acts as a struct for injecting an implementation of UserRepository
// for use in service methods
type userService struct {
	UserRepository repository.IUserRepository
	MailClient     mail.IMailClient
	USVerificationConfig
}

// USConfig will hold repositories that will eventually be injected into
// this service layer
type USConfig struct {
	UserRepository repository.IUserRepository
	MailClient     mail.IMailClient
	USVerificationConfig
}

type USVerificationConfig struct {
	// EmailVerificationTokenSecret string
	// EmailTokenExpirationSecs     int64
	OTPValueFrom      int
	OTPMaxValue       int
	OTPExpirationSecs int64
}

func NewUserService(c *USConfig) IUserService {
	return &userService{
		UserRepository:       c.UserRepository,
		MailClient:           c.MailClient,
		USVerificationConfig: c.USVerificationConfig,
	}
}

// Get retrieves a user based on their uid
func (s *userService) Get(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	return s.UserRepository.FindByID(ctx, uid)
}

func (s *userService) List(ctx context.Context) ([]model.User, error) {
	return s.UserRepository.List(ctx)
}

// HardDelete User
func (s *userService) HardDelete(ctx context.Context, uid uuid.UUID) error {
	return s.UserRepository.HardDeleteUser(ctx, uid)
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.UserRepository.FindByEmail(ctx, email)
}

func (s *userService) ComparePassword(ctx context.Context, email string, password string) (*model.User, error) {
	var uFetched *model.User
	var err error
	if uFetched, err = s.UserRepository.FindByEmail(ctx, email); err != nil {
		apperr, ok := err.(*apperror.Error)
		errstack := apperror.ErrorWrapper(err, helper.TraceCurrentFuncArgs(map[string]string{
			"email": email,
		}))
		if ok {
			return nil, apperror.NewWrapErrorMsg(apperr, errstack)
		}
		return nil, apperror.NewInternalWrap(errstack)
	}

	// verify password
	err = security.ComparePasswords(uFetched.Password, password)
	if err != nil {
		apperr, ok := err.(*apperror.Error)
		errstack := apperror.ErrorWrapper(err, helper.TraceCurrentFuncArgs(map[string]string{
			"email": email,
		}))
		if ok {
			return nil, apperror.NewWrapErrorMsg(apperr, errstack)
		}
		return nil, apperror.NewInternalWrap(errstack)

	}
	// if uFetched.EmailVerifiedAt == nil || uFetched.AccountType == model.TUserUnverified {
	// 	s.regenerateOTPAndSendViaEmail(ctx, *uFetched)
	// 	return nil, &apperror.Error{Type: apperror.TUserUnverified, Message: apperror.UserUnverified}
	// }
	if uFetched.RecordFlag != model.RecActive {
		return nil, apperror.NewAuthorization(apperror.TUserNotActive)
	}

	return uFetched, nil
}

// func (s *userService) FindUserSessionDBByUserID(ctx context.Context, uid uuid.UUID) (*model.UserSession, error) {
// 	return s.UserRepository.FindUserSessionDBByUserID(ctx, uid)
// }

func (s *userService) UpdateUser(ctx context.Context, uid uuid.UUID, uUpdate *model.UserUpdate) error {
	return s.UserRepository.Update(ctx, uid, uUpdate)
}

func (s *userService) StoreSession(ctx context.Context, sess *model.UserSession) error {
	if err := s.UserRepository.StoreSession(ctx, sess); err != nil {
		apperr, ok := err.(*apperror.Error)
		errstack := apperror.ErrorWrapper(err, helper.TraceCurrentFuncArgs(map[string]string{
			"user_id": sess.UserID.String(),
		}))
		if ok {
			return apperror.NewWrapErrorMsg(apperr, errstack)
		}
		return apperror.NewInternalWrap(errstack)
	}
	// s.RedisRepository.SetUserSession(ctx, sess)

	return nil
}

// Logout reaches out to the repository layer to delete all valid tokens for a user
func (s *userService) Logout(ctx context.Context, uid uuid.UUID) error {
	// s.RedisRepository.RemoveUserSessions(ctx, uid)
	if err := s.UserRepository.UpdateUserSession(
		ctx,
		uid,
		model.UserSessionUpdate{Status: model.TUserSessionLogout}); err != nil {
		return err
	}
	return nil
}

func (s *userService) ForgotPasswordUsingEmail(ctx context.Context, email string) error {
	// check if email is exists

	user, err := s.UserRepository.FindByEmail(ctx, email)
	if err != nil {

		apperr, ok := err.(*apperror.Error)
		errstack := apperror.ErrorWrapper(err, helper.TraceCurrentFuncArgs(
			map[string]string{
				"email": email,
			},
		))

		if ok {
			if apperr.Type == apperror.NotFound {
				return &apperror.Error{Type: apperror.TUserNotFound, Message: errstack}
			}
		}
		return apperror.NewInternal()
	}

	// generate a new otp
	otp, err := generateForgotPasswordOTP(s.OTPExpirationSecs, s.OTPValueFrom, s.USVerificationConfig.OTPMaxValue)
	if err != nil {
		errstack := apperror.ErrorWrapper(fmt.Errorf("generate forgot password otp"), helper.TraceCurrentFuncArgs(map[string]string{
			"email": email,
		}))

		return apperror.NewInternalWrap(errstack)
	}
	// prepare a new user otp object

	uOTP := model.UserOTP{
		UserID:         user.ID,
		OTP:            otp.OTP,
		Type:           model.TotpForgotPassword,
		ExpiredAt:      otp.ExpiresAt,
		Status:         model.TOTPCreated,
		DefaultColumns: model.NewDefaultColumn(),
	}

	// store a new user otp
	err = s.UserRepository.UpsertOTP(ctx, user.ID.String(), uOTP)
	if err != nil {
		apperr, ok := err.(*apperror.Error)
		errstack := apperror.ErrorWrapper(err, helper.TraceCurrentFuncArgs(
			map[string]string{
				"email": email,
			},
		))
		if ok {
			return apperror.NewWrapErrorMsg(apperr, errstack)
		}
		return apperror.NewInternalWrap(errstack)
	}

	// send forgot password mail
	err = s.MailClient.SendForgotPasswordOTP(ctx, uOTP.OTP, email)
	if err != nil {
		apperr, ok := err.(*apperror.Error)
		errstack := apperror.ErrorWrapper(err, helper.TraceCurrentFuncArgs(
			map[string]string{
				"email": email,
			},
		))
		if ok {
			return apperror.NewWrapErrorMsg(apperr, errstack)
		}
		return apperror.NewInternalWrap(errstack)
	}

	return nil
}

func (s *userService) UpdatePasswordUsingEmailAndOTP(ctx context.Context, email string, otp string, uUpdate model.UserPasswordUpdate) (*model.User, error) {
	uFetch, err := s.UserRepository.FindByEmail(ctx, email)
	if err != nil {
		apperr, ok := err.(*apperror.Error)
		errstack := apperror.ErrorWrapper(err, helper.TraceCurrentFuncArgs(
			map[string]string{
				"email": email,
			},
		))
		if ok {
			return nil, apperror.NewWrapErrorMsg(apperr, errstack)
		}
		return nil, apperror.NewInternalWrap(errstack)
	}

	err = s.UserRepository.UpdatePassword(ctx, uFetch.ID, uUpdate)
	if err != nil {
		apperr, ok := err.(*apperror.Error)
		errstack := apperror.ErrorWrapper(err, helper.TraceCurrentFuncArgs(
			map[string]string{
				"email": email,
			},
		))
		if ok {
			return nil, apperror.NewWrapErrorMsg(apperr, errstack)
		}
		return nil, apperror.NewInternalWrap(errstack)
	}

	uOTP, err := s.UserRepository.FindOTP(ctx, uFetch.ID, otp, model.TotpForgotPassword)
	if err != nil {
		apperr, ok := err.(*apperror.Error)
		errstack := apperror.ErrorWrapper(err, helper.TraceCurrentFuncArgs(
			map[string]string{
				"email": email,
			},
		))
		if ok {
			return nil, apperror.NewWrapErrorMsg(apperr, errstack)
		}
		return nil, apperror.NewInternalWrap(errstack)
	}

	// tell client if the user otp has been used (soft delete)
	// tell client that otp status must be "VERIFIED"
	if uOTP.RecordFlag == model.RecDeleted || uOTP.Status != model.TOTPVerified {
		return nil, apperror.NewResourceNotFound()
	}

	now := time.Now()
	// if uFetch.AccountType == model.TUserUnverified {
	// 	err = s.UserRepository.Update(ctx, uFetch.ID, &model.UserUpdate{
	// 		AccountType:     model.TUserFree,
	// 		EmailVerifiedAt: &now,
	// 		DefaultColumns: model.DefaultColumns{
	// 			UpdatedAt: &now,
	// 			UpdatedBy: helper.TraceCurrentFunc(),
	// 		},
	// 	})
	// 	if err != nil {
	// 		apperr, ok := err.(*apperror.Error)
	// 		errstack := apperror.ErrorWrapper(err, helper.TraceCurrentFuncArgs(
	// 			map[string]string{
	// 				"email": email,
	// 			},
	// 		))
	// 		if ok {
	// 			return nil, apperror.NewWrapErrorMsg(apperr, errstack)
	// 		}
	// 		return nil, apperror.NewInternalWrap(errstack)
	// 	}
	// }
	err = s.UserRepository.UpdateOTP(ctx, uOTP.ID, model.UserOTPUpdate{
		Status: model.TOTPStatusFinished,
		DefaultColumns: model.DefaultColumns{
			UpdatedAt: &now,
		},
	})
	if err != nil {
		apperr, ok := err.(*apperror.Error)
		errstack := apperror.ErrorWrapper(err, helper.TraceCurrentFuncArgs(
			map[string]string{
				"email": email,
			},
		))
		if ok {
			return nil, apperror.NewWrapErrorMsg(apperr, errstack)
		}
		return nil, apperror.NewInternalWrap(errstack)
	}

	return uFetch, nil
}

func (s *userService) SubmitOTPForgotPassword(ctx context.Context, email string, otp string) error {
	uFetch, err := s.UserRepository.FindByEmail(ctx, email)
	if err != nil {
		apperr, ok := err.(*apperror.Error)
		errstack := apperror.ErrorWrapper(err, helper.TraceCurrentFuncArgs(
			map[string]string{
				"email": email,
			},
		))
		if ok {
			if apperr.Type == apperror.NotFound {
				return &apperror.Error{
					Type:    apperror.EmailNotFound,
					Message: errstack,
				}
			}
			return apperror.NewWrapErrorMsg(apperr, errstack)
		}
		return apperror.NewInternalWrap(errstack)

	}

	uOTP, err := s.UserRepository.FindOTP(ctx, uFetch.ID, otp, model.TotpForgotPassword)
	if err != nil {
		apperr, ok := err.(*apperror.Error)
		errstack := apperror.ErrorWrapper(err, helper.TraceCurrentFuncArgs(
			map[string]string{
				"email": email,
			},
		))
		if ok {
			if apperr.Type == apperror.NotFound {
				return &apperror.Error{
					Type:    apperror.OTPNotFound,
					Message: errstack,
				}
			}
			return apperror.NewWrapErrorMsg(apperr, errstack)
		}
		return apperror.NewInternalWrap(errstack)

	}

	now := time.Now()

	// tell client if the user otp has expired

	if now.UTC().After(uOTP.ExpiredAt) {
		return &apperror.Error{
			Type: apperror.OTPExpired,
			Message: apperror.ErrorWrapper(fmt.Errorf("forgot password otp expired"), helper.TraceCurrentFuncArgs(
				map[string]string{
					"email": email,
				},
			)),
		}
	}
	// tell client if the user otp has been used (soft delete)
	// tell client that otp status must in "CREATED"
	if uOTP.RecordFlag == model.RecDeleted || uOTP.Status != model.TOTPCreated {
		return &apperror.Error{
			Type: apperror.OTPNotFound,
			Message: apperror.ErrorWrapper(fmt.Errorf("forgot password otp deleted/not created"), helper.TraceCurrentFuncArgs(
				map[string]string{
					"email": email,
				},
			)),
		}
	}

	return s.UserRepository.UpdateOTP(
		ctx,
		uOTP.ID,
		model.UserOTPUpdate{
			Status: model.TOTPVerified, DefaultColumns: model.DefaultColumns{
				UpdatedAt: &now,
			},
		},
	)
}

func (s *userService) UpdateUserSession(ctx context.Context, uid uuid.UUID, uSessUPD model.UserSessionUpdate) error {
	if err := s.UserRepository.UpdateUserSession(ctx, uid, uSessUPD); err != nil {
		apperr, ok := err.(*apperror.Error)
		errstack := apperror.ErrorWrapper(err, helper.TraceCurrentFuncArgs(map[string]string{
			"user_id": uid.String(),
		}))
		if ok {
			return apperror.NewWrapErrorMsg(apperr, errstack)
		}
		return apperror.NewInternalWrap(errstack)
	}
	uSess, err := uSessUPD.ToUserSession(uid)
	if err != nil {
		return err
	}
	uSess.UserID = uid
	// s.RedisRepository.SetUserSession(ctx, uSess)
	return nil
}

// func (s *userService) ResendEmailVerification(ctx context.Context, u model.User) error {
// 	user, err := s.UserRepository.FindByEmail(ctx, u.Email)
// 	// check whether user account_type is unverified
// 	if err != nil {
// 		return err
// 	}

// 	if user.AccountType != model.TUserUnverified {

// 		errstack := apperror.ErrorWrapper(fmt.Errorf("user is not in unverified state"), helper.TraceCurrentFuncArgs(u))

// 		return apperror.NewConflictMsg(errstack)
// 	}

// 	otp, err := generateEmailVerficationOTP(s.EmailTokenExpirationSecs, s.OTPValueFrom, s.OTPMaxValue)

// 	emailOTP := model.UserOTP{
// 		UserID:         user.ID,
// 		OTP:            otp.OTP,
// 		Type:           model.TotpRegistVerifyEmail,
// 		ExpiredAt:      otp.ExpiresAt,
// 		Status:         model.TOTPCreated,
// 		DefaultColumns: model.NewDefaultColumn(),
// 	}

// 	err = s.UserRepository.UpsertOTP(ctx, user.ID.String(), emailOTP)
// 	if err != nil {
// 		apperr, ok := err.(*apperror.Error)
// 		errstack := apperror.ErrorWrapper(err, helper.TraceCurrentFuncArgs(
// 			map[string]string{
// 				"user_id": u.ID.String(),
// 			},
// 		))
// 		if ok {
// 			return apperror.NewWrapErrorMsg(apperr, errstack)
// 		}
// 		return apperror.NewInternalWrap(errstack)
// 	}

// 	err = s.MailClient.SendEmailVerification(ctx, otp.OTP, u.Email)
// 	if err != nil {
// 		apperr, ok := err.(*apperror.Error)
// 		errstack := apperror.ErrorWrapper(err, helper.TraceCurrentFuncArgs(
// 			map[string]string{
// 				"user_id": u.ID.String(),
// 			},
// 		))
// 		if ok {
// 			return apperror.NewWrapErrorMsg(apperr, errstack)
// 		}
// 		return apperror.NewInternalWrap(errstack)
// 	}
// 	return nil
// }

// Register reaches our to a UserRepository to verify the
// email address is avaliable and signs up the user if this is the case
func (s *userService) Register(ctx context.Context, u model.User) (*model.User, *model.OTPData, error) {
	// check if record is available
	userExists, err := s.UserRepository.EmailExists(ctx, u.Email)
	if err != nil {
		apperr, ok := err.(*apperror.Error)
		errstack := apperror.ErrorWrapper(err, helper.TraceCurrentFuncArgs(u))
		if ok {
			return nil, nil, apperror.NewWrapErrorMsg(apperr, errstack)
		}
		return nil, nil, apperror.NewInternalWrap(errstack)
	}
	if userExists {
		return nil, nil, apperror.NewConflictSimple()
	}
	// set user to unverified

	otp, err := generateEmailVerficationOTP(1000, s.OTPValueFrom, s.OTPMaxValue)
	if err != nil {
		errstack := apperror.ErrorWrapper(fmt.Errorf("generateEmailVerficationOTP"), helper.TraceCurrentFuncArgs(u))
		return nil, nil, apperror.NewInternalWrap(errstack)
	}
	// u.AccTypeUnverified()
	u.CreatedBy = helper.TraceCurrentFunc()
	userCreated, err := s.UserRepository.Create(ctx, u)
	if err != nil {
		apperr, ok := err.(*apperror.Error)
		errstack := apperror.ErrorWrapper(err, helper.TraceCurrentFuncArgs(u))
		if ok {
			return nil, nil, apperror.NewWrapErrorMsg(apperr, errstack)
		}
		return nil, nil, apperror.NewInternalWrap(errstack)
	}
	emailOTP := model.UserOTP{
		UserID:         userCreated.ID,
		OTP:            otp.OTP,
		Type:           model.TotpRegistVerifyEmail,
		ExpiredAt:      otp.ExpiresAt,
		Status:         model.TOTPCreated,
		DefaultColumns: model.NewDefaultColumn(),
	}

	err = s.UserRepository.UpsertOTP(ctx, userCreated.ID.String(), emailOTP)
	if err != nil {

		apperr, ok := err.(*apperror.Error)
		errstack := apperror.ErrorWrapper(err, helper.TraceCurrentFuncArgs(u))
		if ok {
			return nil, nil, apperror.NewWrapErrorMsg(apperr, errstack)
		}
		return nil, nil, apperror.NewInternalWrap(errstack)
	}

	// If we get around to adding events, we'd publish it here
	// err := s.EventsBroker PublishUserUpdated(u, true)
	return userCreated, otp, nil
}

// func (s *userService) StoreFromSocialLogin(ctx context.Context, socialName, nickname, socialEmail, gender string) (*model.User, error) {
// 	var (
// 		finalName  string
// 		userExists bool = true
// 	)
// 	finalGender, ok := model.NewGender(gender)
// 	if !ok {
// 		errstack := apperror.ErrorWrapper(fmt.Errorf("NewGender"), helper.TraceCurrentFuncArgs(map[string]string{
// 			"social_name":  socialName,
// 			"nick_name":    nickname,
// 			"social_email": socialEmail,
// 			"gender":       gender,
// 		}))
// 		return nil, apperror.NewBadRequest(errstack)
// 	}
// 	// BEGIN: user exists
// 	user, err := s.UserRepository.FindByEmail(ctx, socialEmail)
// 	if err != nil {
// 		apperr, ok := err.(*apperror.Error)
// 		errstack := apperror.ErrorWrapper(fmt.Errorf("NewGender"), helper.TraceCurrentFuncArgs(map[string]string{
// 			"social_name":  socialName,
// 			"nick_name":    nickname,
// 			"social_email": socialEmail,
// 			"gender":       gender,
// 		}))

// 		if ok {
// 			// user not found
// 			if apperr.Type == apperror.NotFound {
// 				userExists = false
// 			} else {
// 				return nil, apperror.NewWrapErrorMsg(apperr, errstack)
// 			}
// 		} else {
// 			return nil, apperror.NewInternalWrap(errstack)
// 		}
// 	}

// 	// user exists
// 	if !userExists {

// 		finalName = socialName
// 		if len(finalName) == 0 {
// 			finalName = nickname
// 		}

// 		now := time.Now()
// 		u := model.User{
// 			AccountType:     model.TUserFree,
// 			Name:            finalName,
// 			Gender:          finalGender,
// 			Nickname:        nickname,
// 			Email:           socialEmail,
// 			Password:        "",
// 			EmailVerifiedAt: &now,
// 			DefaultColumns:  model.NewDefaultColumn(),
// 		}

// 		user, err := s.UserRepository.Create(ctx, u)
// 		if err != nil {
// 			apperr, ok := err.(*apperror.Error)
// 			errstack := apperror.ErrorWrapper(err, helper.TraceCurrentFuncArgs(map[string]string{
// 				"social_name":  socialName,
// 				"nick_name":    nickname,
// 				"social_email": socialEmail,
// 				"gender":       gender,
// 			}))
// 			if ok {
// 				return nil, apperror.NewWrapErrorMsg(apperr, errstack)
// 			}
// 			return nil, apperror.NewInternalWrap(errstack)
// 		}
// 		// s.MailClient.SendOnboardingGreeting(ctx, u.Email)
// 		return user, nil
// 	}

// 	if len(user.Name) == 0 && len(socialName) == 0 {
// 		finalName = nickname
// 	} else if len(socialName) != 0 {
// 		finalName = socialName
// 	}

// 	now := time.Now()
// 	updateUser := model.UserUpdate{
// 		Name: finalName,
// 		DefaultColumns: model.DefaultColumns{
// 			UpdatedAt: &now,
// 			UpdatedBy: helper.TraceCurrentFunc(),
// 		},
// 	}

// 	// make user type: verified
// 	if user.AccountType == model.TUserUnverified {
// 		updateUser.AccountType = model.TUserFree
// 		updateUser.EmailVerifiedAt = &now
// 	}
// 	if err := s.UserRepository.Update(ctx, user.ID, &updateUser); err != nil {
// 		apperr, ok := err.(*apperror.Error)
// 		errstack := apperror.ErrorWrapper(err, helper.TraceCurrentFuncArgs(map[string]string{
// 			"social_name":  socialName,
// 			"nick_name":    nickname,
// 			"social_email": socialEmail,
// 			"gender":       gender,
// 		}))
// 		if ok {
// 			return nil, apperror.NewWrapErrorMsg(apperr, errstack)
// 		}
// 		return nil, apperror.NewInternalWrap(errstack)
// 	}
// 	return user, nil
// }

// func (s *userService) SocialLoginUpdateOrStore(ctx context.Context, socLogID string, socialLogin model.SocialLogin) (*model.SocialLogin, error) {
// 	socLoginFetch, err := s.UserRepository.GetSocialLoginBySocialUserID(ctx, socLogID)
// 	if err != nil {
// 		apperr, ok := err.(*apperror.Error)
// 		errstack := apperror.ErrorWrapper(err, helper.TraceCurrentFuncArgs(map[string]string{
// 			"social_login_id": socLogID,
// 			"user_id":         socialLogin.UserID.String(),
// 		}))

// 		if ok {
// 			if apperr.Type == apperror.NotFound {
// 				return s.UserRepository.SocialLoginStore(ctx, socialLogin)
// 			}
// 			return nil, apperror.NewWrapErrorMsg(apperr, errstack)
// 		}
// 		return nil, apperror.NewInternalWrap(errstack)
// 	}
// 	now := time.Now()
// 	socLoginUpdate := model.SocialLoginUpdate{
// 		SocialToken: socialLogin.SocialToken,
// 		DefaultColumns: model.DefaultColumns{
// 			UpdatedAt: &now,
// 			UpdatedBy: helper.TraceCurrentFunc(),
// 		},
// 	}
// 	err = s.UserRepository.SocialLoginUpdate(ctx, socLoginFetch.ID, socLoginUpdate)
// 	if err != nil {
// 		apperr, ok := err.(*apperror.Error)
// 		errstack := apperror.ErrorWrapper(err, helper.TraceCurrentFuncArgs(map[string]string{
// 			"social_login_id": socLogID,
// 			"user_id":         socialLogin.UserID.String(),
// 		}))
// 		if ok {
// 			return nil, apperror.NewWrapErrorMsg(apperr, errstack)
// 		}
// 		return nil, apperror.NewInternalWrap(errstack)
// 	}
// 	return nil, nil
// }

// Signin: 1. check if the user exists
// 2. compare the supplied password with the provided password
// 3. if a valid email/password combination is provided, return the fetched User
// available user fields

// func (s *userService) regenerateOTPAndSendViaEmail(ctx context.Context, userCreated model.User) error {
// 	otp, err := generateEmailVerficationOTP(s.EmailTokenExpirationSecs, s.OTPValueFrom, s.OTPMaxValue)
// 	if err != nil {
// 		errstack := apperror.ErrorWrapper(fmt.Errorf("generate email otp"), helper.TraceCurrentFuncArgs(userCreated.Email))
// 		return apperror.NewInternalWrap(errstack)
// 	}

// 	emailOTP := model.UserOTP{
// 		UserID:         userCreated.ID,
// 		OTP:            otp.OTP,
// 		Type:           model.TotpRegistVerifyEmail,
// 		ExpiredAt:      otp.ExpiresAt,
// 		Status:         model.TOTPCreated,
// 		DefaultColumns: model.NewDefaultColumn(),
// 	}

// 	err = s.UserRepository.UpsertOTP(ctx, userCreated.ID.String(), emailOTP)
// 	if err != nil {
// 		apperr, ok := err.(*apperror.Error)
// 		errstack := apperror.ErrorWrapper(err, helper.TraceCurrentFuncArgs(map[string]string{
// 			"email": userCreated.Email,
// 		}))
// 		if ok {
// 			return apperror.NewWrapErrorMsg(apperr, errstack)
// 		}
// 		return apperror.NewInternalWrap(errstack)
// 	}

// 	err = s.MailClient.SendEmailVerification(ctx, otp.OTP, userCreated.Email)
// 	if err != nil {
// 		apperr, ok := err.(*apperror.Error)
// 		errstack := apperror.ErrorWrapper(err, helper.TraceCurrentFuncArgs(map[string]string{
// 			"email": userCreated.Email,
// 		}))
// 		if ok {
// 			return apperror.NewWrapErrorMsg(apperr, errstack)
// 		}
// 		return apperror.NewInternalWrap(errstack)
// 	}

// 	return nil
// }

// func (s *userService) RefreshTokenExists(ctx context.Context, uid uuid.UUID, refreshTokenID uuid.UUID) (bool, error) {
// 	exists, err := s.RedisRepository.RefreshTokenExists(ctx, uid, refreshTokenID)
// 	// if redis token exists, then no skip the db repo checking
// 	if exists {
// 		return true, nil
// 	}
// }

// 	exists, err = s.UserRepository.RefreshTokenExists(ctx, uid, refreshTokenID)
// 	if err != nil {
// 		apperr, ok := err.(*apperror.Error)
// 		errstack := apperror.ErrorWrapper(err, helper.TraceCurrentFuncArgs(map[string]string{
// 			"user_id":          uid.String(),
// 			"refresh_token_id": refreshTokenID.String(),
// 		}))
// 		if ok {
// 			return false, apperror.NewWrapErrorMsg(apperr, errstack)
// 		}
// 		return false, apperror.NewInternalWrap(errstack)
// 	}

// 	if exists {
// 		return true, nil
// 	}
// 	return false, nil
// }

// func (s *userService) RedisRemoveSession(ctx context.Context, uid uuid.UUID) error {
// 	return s.RedisRepository.RemoveUserSessions(ctx, uid)
// }

// func (s *userService) RedisGetSession(ctx context.Context, uid uuid.UUID, refTokenID uuid.UUID) (*model.UserSession, error) {
// 	return s.RedisRepository.GetUserSession(ctx, uid, refTokenID)
// }

// func (s *userService) VerifyEmail(ctx context.Context, email, otp string) (*model.User, error) {
// 	uFetch, err := s.UserRepository.FindByEmail(ctx, email)
// 	if err != nil {

// 		apperr, ok := err.(*apperror.Error)
// 		errstack := apperror.ErrorWrapper(err, helper.TraceCurrentFuncArgs(map[string]string{
// 			"email": email,
// 		}))
// 		if ok {
// 			if apperr.Type == apperror.NotFound {
// 				return nil, &apperror.Error{
// 					Type:    apperror.EmailNotFound,
// 					Message: errstack,
// 				}
// 			}
// 			return nil, apperror.NewWrapErrorMsg(apperr, errstack)
// 		}
// 		return nil, apperror.NewInternalWrap(errstack)

// 	}

// 	if uFetch.AccountType != model.TUserUnverified {
// 		return nil, apperror.NewResourceNotFound()
// 	}

// 	otpFetch, err := s.UserRepository.FindOTP(ctx, uFetch.ID, otp, model.TotpRegistVerifyEmail)
// 	if err != nil {
// 		apperr, ok := err.(*apperror.Error)
// 		errstack := apperror.ErrorWrapper(err, helper.TraceCurrentFuncArgs(map[string]string{
// 			"email": email,
// 		}))
// 		if ok {
// 			if apperr.Type == apperror.NotFound {
// 				return nil, &apperror.Error{
// 					Type:    apperror.OTPNotFound,
// 					Message: errstack,
// 				}
// 			}
// 			return nil, apperror.NewWrapErrorMsg(apperr, errstack)
// 		}
// 		return nil, apperror.NewInternalWrap(errstack)
// 	}

// 	now := time.Now()
// 	if now.UTC().After(otpFetch.ExpiredAt) {
// 		errstack := apperror.ErrorWrapper(fmt.Errorf("email otp has expired"), helper.TraceCurrentFuncArgs(map[string]string{
// 			"email": email,
// 		}))

// 		return nil, &apperror.Error{
// 			Type:    apperror.OTPExpired,
// 			Message: errstack,
// 		}
// 	}

// 	if otpFetch.RecordFlag == model.RecDeleted || otpFetch.Status != model.TOTPCreated {
// 		errstack := apperror.ErrorWrapper(fmt.Errorf("invalid email otp"), helper.TraceCurrentFuncArgs(map[string]string{
// 			"email": email,
// 		}))

// 		return nil, &apperror.Error{
// 			Type:    apperror.OTPNotFound,
// 			Message: errstack,
// 		}
// 	}

// 	err = s.UserRepository.Update(ctx, uFetch.ID, &model.UserUpdate{
// 		AccountType:     model.TUserFree,
// 		EmailVerifiedAt: &now,
// 		DefaultColumns: model.DefaultColumns{
// 			UpdatedAt: &now,
// 			UpdatedBy: helper.TraceCurrentFunc(),
// 		},
// 	})
// 	if err != nil {
// 		apperr, ok := err.(*apperror.Error)
// 		errstack := apperror.ErrorWrapper(err, helper.TraceCurrentFuncArgs(
// 			map[string]string{
// 				"email": email,
// 			},
// 		))
// 		if ok {
// 			return nil, apperror.NewWrapErrorMsg(apperr, errstack)
// 		}
// 		return nil, apperror.NewInternalWrap(errstack)
// 	}

// 	err = s.UserRepository.UpdateOTP(ctx, otpFetch.ID,
// 		model.UserOTPUpdate{
// 			Status: model.TOTPStatusFinished,
// 			DefaultColumns: model.DefaultColumns{
// 				UpdatedAt: &now,
// 			},
// 		})
// 	if err != nil {
// 		apperr, ok := err.(*apperror.Error)
// 		errstack := apperror.ErrorWrapper(err, helper.TraceCurrentFuncArgs(
// 			map[string]string{
// 				"email": email,
// 			},
// 		))
// 		if ok {
// 			return nil, apperror.NewWrapErrorMsg(apperr, errstack)
// 		}
// 		return nil, apperror.NewInternalWrap(errstack)
// 	}
// 	// s.MailClient.SendOnboardingGreeting(ctx, email)

// 	return uFetch, nil
// }

// func (s *userService) IsSubscribing(ctx context.Context, userID uuid.UUID) (bool, error) {
// 	return s.UserRepository.IsSubscribing(ctx, userID)
// }
