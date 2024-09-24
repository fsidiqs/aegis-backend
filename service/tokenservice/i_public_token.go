package tokenservice

import (
	"context"

	"github.com/fsidiqs/aegis-backend/model"
)

type IPublicTokenService interface {
	NewPairFromDeviceID(ctx context.Context, deviceID, notifToken string) (*model.PublicTokenData, error)
	ValidatePublicToken(publicTokenStr string) (*model.ValidatedPublicToken, error)
}
