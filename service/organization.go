package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/fsidiqs/aegis-backend/helper"
	"github.com/fsidiqs/aegis-backend/model"
	"github.com/fsidiqs/aegis-backend/model/apperror"
	"github.com/fsidiqs/aegis-backend/repository"
)

// UserService acts as a struct for injecting an implementation of UserRepository
// for use in service methods
type organizationService struct {
	OrganizationRepository repository.IOrganizationRepository
}

// USConfig will hold repositories that will eventually be injected into
// this service layer
type OSConfig struct {
	OrganizationRepository repository.IOrganizationRepository
}

func NewOrganizationService(c *OSConfig) IOrganizationService {
	return &organizationService{
		OrganizationRepository: c.OrganizationRepository,
	}
}

// Get retrieves a user based on their uid
func (s *organizationService) Get(ctx context.Context, uid uuid.UUID) (*model.Organization, error) {
	return s.OrganizationRepository.FindByID(ctx, uid)
}

func (s *organizationService) List(ctx context.Context) ([]model.Organization, error) {
	return s.OrganizationRepository.List(ctx)
}

// HardDelete User
func (s *organizationService) HardDelete(ctx context.Context, uid uuid.UUID) error {
	return s.OrganizationRepository.HardDelete(ctx, uid)
}

func (s *organizationService) Create(ctx context.Context, u model.Organization, userId uuid.UUID) (*model.Organization, error) {
	// check if record is available
	// set user to unverified
	u.CreatorID = userId
	userCreated, err := s.OrganizationRepository.Create(ctx, u, userId)
	if err != nil {
		apperr, ok := err.(*apperror.Error)
		errstack := apperror.ErrorWrapper(err, helper.TraceCurrentFuncArgs(u))
		if ok {
			return nil, apperror.NewWrapErrorMsg(apperr, errstack)
		}
		return nil, apperror.NewInternalWrap(errstack)
	}
	return userCreated, nil
}

func (s *organizationService) Update(ctx context.Context, uid uuid.UUID, orgUpdate *model.Organization) error {
	return s.OrganizationRepository.Update(ctx, uid, orgUpdate)
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
