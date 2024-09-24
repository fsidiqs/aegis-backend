package service

import (
	"context"
	"crypto/rsa"
	"fmt"

	"firebase.google.com/go/auth"
	"github.com/google/uuid"

	"github.com/fsidiqs/aegis-backend/helper"
	"github.com/fsidiqs/aegis-backend/model"

	"github.com/fsidiqs/aegis-backend/model/apperror"
)

// TokenService used for injecting an implementation of RedisRepository
// for use in service methdos along with keys and secrets for
// signing JWTs
type tokenServiceImpl struct {
	PrivKey               *rsa.PrivateKey
	PubKey                *rsa.PublicKey
	RefreshSecret         string
	IDExpirationSecs      int64
	RefreshExpirationSecs int64
}

// TSConfig will hold repositories that will eventually be injected into
// this service layer
type TSConfig struct {
	PrivKey               *rsa.PrivateKey
	PubKey                *rsa.PublicKey
	RefreshSecret         string
	IDExpirationSecs      int64
	RefreshExpirationSecs int64
}

func NewTokenService(c *TSConfig) ITokenService {
	return &tokenServiceImpl{
		PrivKey:               c.PrivKey,
		PubKey:                c.PubKey,
		RefreshSecret:         c.RefreshSecret,
		IDExpirationSecs:      c.IDExpirationSecs,
		RefreshExpirationSecs: c.RefreshExpirationSecs,
	}
}

// NewPairFromUser creates freash id and refresh tokens for the current user
// If a previous token is included, the previous token is removed from
// the tokens repository
// 1. RedisRepository.DeleteRefreshToken will check whether the previousTokenID exists or not, not exist meanings invalid and return the error
// 2. If the previous token id exists, then we generate the id token signed with privatekey
// 3. we generate the refreshtoken, signed with refresh secret string
// 4. Then we set the token repository with the generaetd refreshToken
func (s *tokenServiceImpl) NewPairFromUser(ctx context.Context, u *model.User) (*model.TokenData, error) {
	refreshToken, err := generateRefreshToken(u.ID, s.RefreshSecret, s.RefreshExpirationSecs)
	if err != nil {
		errmmsg := fmt.Sprintf("%s", helper.TraceCurrentFunc())
		return nil, apperror.NewInternalWrap(errmmsg)
	}

	authToken, err := generateAuthToken(u, refreshToken.ID, s.PrivKey, s.IDExpirationSecs)
	if err != nil {
		errmmsg := fmt.Sprintf("%s", helper.TraceCurrentFunc())
		return nil, apperror.NewInternalWrap(errmmsg)
	}

	return &model.TokenData{
		AuthToken:    model.AuthToken{SS: authToken},
		RefreshToken: model.RefreshToken{SS: refreshToken.SS, ID: refreshToken.ID, UID: u.ID, ExpiresIn: refreshToken.ExpiresIn},
	}, nil
}

// ValidateAuthToken validates the id token jwt string
// It returns the user extract from the AuthTokenCustomClaims
func (s *tokenServiceImpl) ValidateAuthToken(tokenString string) (*model.ValidatedAuthData, error) {
	claims, err := validateAuthToken(tokenString, s.PubKey) // uses public RSA key
	// We'll just return unauthorized error in all instances of failing to verify user
	if err != nil {
		errMsg := fmt.Sprintf("%s", helper.TraceCurrentFunc())
		return nil, apperror.NewAuthorization(errMsg)
	}

	ret := &model.ValidatedAuthData{
		User:           claims.User,
		RefreshTokenID: claims.RefreshTokenID,
		Provider:       claims.Provider,
	}
	return ret, nil
}

// ValidateRefreshToken checks to make sure the JWT Provided by a string is valid
// and returns a RefreshToken if valid
func (s *tokenServiceImpl) ValidateRefreshToken(tokenString string) (*model.RefreshToken, error) {
	// validate actual JWT with a secret string
	claims, err := validateRefreshToken(tokenString, s.RefreshSecret)
	// We'll just return unauthorized error in all instances of failing to verify user
	if err != nil {
		errMsg := fmt.Sprintf("%s", helper.TraceCurrentFunc())
		return nil, apperror.NewAuthorization(errMsg)
	}

	// Standard claims store ID as a string. So we parse claims.Id as UUID

	tokenUUID, err := uuid.Parse(claims.Id)
	if err != nil {
		errMsg := fmt.Sprintf("%s", helper.TraceCurrentFunc())
		return nil, apperror.NewAuthorization(errMsg)
	}

	return &model.RefreshToken{
		SS:  tokenString,
		ID:  tokenUUID,
		UID: claims.UID,
	}, nil
}

func (s *tokenServiceImpl) FirebaseAuthTokenNewSocialLogin(token *auth.Token) (*model.SocialLoginReq, error) {
	socL := &model.SocialLoginReq{
		Name:         fmt.Sprintf("%s", token.Claims["name"]),
		Email:        fmt.Sprintf("%s", token.Claims["email"]),
		SocialUserID: token.UID,
		Provider:     token.Firebase.SignInProvider,
	}
	return socL, nil
}
