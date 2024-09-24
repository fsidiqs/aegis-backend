package service

import (
	"context"

	"firebase.google.com/go/auth"
	"github.com/fsidiqs/aegis-backend/model"
	"github.com/google/uuid"
)

// UserService defines methods the handler layer expects
// any service it interacts with to implement
type IUserService interface {
	Get(ctx context.Context, uid uuid.UUID) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Register(ctx context.Context, u model.User) (*model.User, *model.OTPData, error)
	// StoreFromSocialLogin(ctx context.Context, name, nickname, email, gender string) (*model.User, error)
	ComparePassword(ctx context.Context, email string, password string) (*model.User, error)
	// Logout(ctx context.Context, uid uuid.UUID) error
	UpdateUser(ctx context.Context, uid uuid.UUID, uReq *model.UserUpdate) error

	// SocialLoginUpdateOrStore(ctx context.Context, socLoginUserID string, socialLogin model.SocialLogin) (*model.SocialLogin, error)

	// FindUserSessionDBByUserID(ctx context.Context, uid uuid.UUID) (*model.UserSession, error)
	StoreSession(ctx context.Context, s *model.UserSession) error
	// RefreshTokenExists(ctx context.Context, uid uuid.UUID, refreshTokenid uuid.UUID) (bool, error)
	UpdateUserSession(ctx context.Context, uid uuid.UUID, uSessUPD model.UserSessionUpdate) error

	// VerifyEmail(ctx context.Context, email, otp string) (*model.User, error)
	// ResendEmailVerification(ctx context.Context, u model.User) error
	ForgotPasswordUsingEmail(ctx context.Context, email string) error
	SubmitOTPForgotPassword(ctx context.Context, email string, otp string) error
	UpdatePasswordUsingEmailAndOTP(ctx context.Context, email string, otp string, uUpdate model.UserPasswordUpdate) (*model.User, error)

	// RedisGetSession(ctx context.Context, uid uuid.UUID, refTokenID uuid.UUID) (*model.UserSession, error)
	// RedisRemoveSession(ctx context.Context, uid uuid.UUID) error

	// IsSubscribing(ctx context.Context, userID uuid.UUID) (bool, error)
}

// TokenService defines methods the handler layer expects to interact
// with in regards to producing JWTs as string

type ITokenService interface {
	NewPairFromUser(ctx context.Context, u *model.User) (*model.TokenData, error)
	ValidateAuthToken(tokenString string) (*model.ValidatedAuthData, error)
	ValidateRefreshToken(refreshTokenString string) (*model.RefreshToken, error)
	FirebaseAuthTokenNewSocialLogin(token *auth.Token) (*model.SocialLoginReq, error)
}
