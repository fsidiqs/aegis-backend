package tokenservice

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fsidiqs/aegis-backend/helper"
	"github.com/fsidiqs/aegis-backend/model/apperror"
)

// PUBLIC TOKEN

type publicTokenCustomClaims struct {
	DeviceID           string `json:"device_id"`
	NotificactionToken string `json:"notification_token"`
	jwt.StandardClaims
}

func generatePublicToken(deviceID, notifToken, secret string, exp int64) (string, error) {
	unixTime := time.Now().Unix()
	tokenExp := unixTime + exp

	claims := publicTokenCustomClaims{
		DeviceID:           deviceID,
		NotificactionToken: notifToken,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "aegis",
			IssuedAt:  unixTime,
			ExpiresAt: tokenExp,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(secret))
	if err != nil {
		errMsg := fmt.Sprintf("%s", helper.TraceCurrentFuncArgs(
			map[string]string{
				"deviceID":   deviceID,
				"notifToken": notifToken,
				"secret":     secret,
			},
		))
		return "", apperror.NewBadRequest(errMsg)
	}
	return ss, nil
}

func validatePublicToken(tokenString, secret string) (*publicTokenCustomClaims, error) {
	claims := &publicTokenCustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, &apperror.Error{Message: apperror.InvalidPublicToken, Type: apperror.Authorization}
	}

	claims, ok := token.Claims.(*publicTokenCustomClaims)

	if !ok {
		return nil, &apperror.Error{Message: apperror.ErrClaimParse, Type: apperror.Authorization}
	}

	return claims, nil
}
