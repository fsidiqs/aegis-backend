package tokenservice

import (
	"context"
	"log"

	helper "github.com/fsidiqs/aegis-backend/helper"
	model "github.com/fsidiqs/aegis-backend/model"
	"github.com/fsidiqs/aegis-backend/model/apperror"
)

type publicTokenServiceImpl struct {
	Secret               string
	SecretExpirationSecs int64
}

type PublicTSConfig struct {
	Secret               string
	SecretExpirationSecs int64
}

func NewPublicTokenService(c *PublicTSConfig) IPublicTokenService {
	return &publicTokenServiceImpl{
		Secret:               c.Secret,
		SecretExpirationSecs: c.SecretExpirationSecs,
	}
}

func (s *publicTokenServiceImpl) NewPairFromDeviceID(ctx context.Context, dvcID, notifToken string) (*model.PublicTokenData, error) {
	publicToken, err := generatePublicToken(dvcID, notifToken, s.Secret, s.SecretExpirationSecs)
	if err != nil {
		errMsg := apperror.ErrorWrapper(err, helper.TraceCurrentFuncArgs(
			map[string]string{
				"deviceID":   dvcID,
				"notifToken": notifToken,
			}))
		return nil, apperror.NewBadRequest(errMsg)
	}

	return &model.PublicTokenData{
		PublicToken: model.PublicToken{SS: publicToken},
	}, nil
}

func (s *publicTokenServiceImpl) ValidatePublicToken(publicTokenStr string) (*model.ValidatedPublicToken, error) {
	claims, err := validatePublicToken(publicTokenStr, s.Secret)
	if err != nil {
		log.Printf("Error %v err:%v", helper.TraceCurrentFunc(), err)
		return nil, err
	}

	return &model.ValidatedPublicToken{
		DeviceID:          claims.DeviceID,
		NotificationToken: claims.NotificactionToken,
	}, nil
}
