package repository

import (
	"context"

	"github.com/fsidiqs/aegis-backend/model"
	"github.com/google/uuid"
)

// UserRepository defines methods the service layer expects
// any repository it interacts with to implement
type IUserRepository interface {
	Create(ctx context.Context, u model.User) (*model.User, error)
	Update(ctx context.Context, uid uuid.UUID, u *model.UserUpdate) error
	FindByID(ctx context.Context, uid uuid.UUID) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	// FindByEmailAndPhone(ctx context.Context, email string, phone string) (*model.User, error)
	UpdatePassword(ctx context.Context, uid uuid.UUID, u model.UserPasswordUpdate) error
	EmailExists(ctx context.Context, email string) (bool, error)

	// RefreshTokenExists(ctx context.Context, uid uuid.UUID, refreshTokenId uuid.UUID) (bool, error)
	StoreSession(ctx context.Context, s *model.UserSession) error
	// FindUserSessionDBByUserID(ctx context.Context, uid uuid.UUID) (*model.UserSession, error)
	RemoveSession(ctx context.Context, userId uuid.UUID, prevTokenId uuid.UUID) error
	// RemoveUserSessions(ctx context.Context, userId uuid.UUID) error
	UpdateUserSession(ctx context.Context, userID uuid.UUID, uSessUpdate model.UserSessionUpdate) error
	// // Forgot Password OTP
	FindOTP(ctx context.Context, userId uuid.UUID, otp string, otpType model.TOTP) (*model.UserOTP, error)
	UpsertOTP(ctx context.Context, userID string, uOTP model.UserOTP) error
	UpdateOTP(ctx context.Context, uOTPID uuid.UUID, userOTP model.UserOTPUpdate) error

	// IsSubscribing(ctx context.Context, userID uuid.UUID) (bool, error)
}
