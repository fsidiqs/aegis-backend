package service

import (
	"crypto/rsa"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fsidiqs/aegis-backend/helper"
	"github.com/fsidiqs/aegis-backend/model"
	"github.com/fsidiqs/aegis-backend/model/apperror"
	"github.com/google/uuid"
)

// AuthTokenCUstomCLaims holds structure of jwt claims of authToken
type authTokenCustomClaims struct {
	User           *model.User `json:"user"`
	RefreshTokenID string      `json:"refresh_token_id"`
	Provider       string      `json:"auth_provider"`
	jwt.StandardClaims
}

func generateAuthToken(u *model.User, refreshTokenID uuid.UUID, key *rsa.PrivateKey, exp int64) (string, error) {
	unixTime := time.Now().Unix()
	tokenExp := unixTime + exp

	claims := authTokenCustomClaims{
		User:           u,
		RefreshTokenID: refreshTokenID.String(),
		Provider:       "aegis",

		StandardClaims: jwt.StandardClaims{
			Issuer:    "aegis",
			IssuedAt:  unixTime,
			ExpiresAt: tokenExp,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	ss, err := token.SignedString(key)
	if err != nil {
		errMsg := fmt.Sprintf("%s:err:%v", helper.TraceCurrentFunc(), err)
		return "", apperror.NewInternalWrap(errMsg)
	}
	return ss, nil
}

// refreshTokenCustomClaims holds the payload of a refresh token
// This can be used to extract user id for subsequent
// application operations (IE, fetch user in Redis)
type refreshTokenCustomClaims struct {
	UID uuid.UUID `json:"uid"`
	jwt.StandardClaims
}

// generateRefreshToken creates a refresh token
// The refresh token stores the signed string, the user's ID (string), and expiry time
func generateRefreshToken(uid uuid.UUID, key string, exp int64) (*model.RefreshToken, error) {
	currentTime := time.Now()
	tokenExp := currentTime.Add(time.Duration(exp) * time.Second)

	tokenID, err := uuid.NewRandom() // v4 uuid in the google uuid lib
	if err != nil {
		errMsg := fmt.Sprintf("%s:err:%v", helper.TraceCurrentFunc(), err)
		return nil, apperror.NewInternalWrap(errMsg)
	}

	claims := refreshTokenCustomClaims{
		UID: uid,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  currentTime.Unix(),
			ExpiresAt: tokenExp.Unix(),
			Id:        tokenID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(key))
	if err != nil {
		errMsg := fmt.Sprintf("%s:err:%v", helper.TraceCurrentFunc(), err)
		return nil, apperror.NewInternalWrap(errMsg)
	}
	return &model.RefreshToken{
		SS:        ss,
		ID:        tokenID,
		UID:       uid,
		ExpiresIn: tokenExp.Sub(currentTime),
	}, nil
}

// validateAuthToken returns the token's claims if the token is valid
func validateAuthToken(tokenString string, key *rsa.PublicKey) (*authTokenCustomClaims, error) {
	claims := &authTokenCustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	// For now we'll just return the error and handle logging in service level
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("ID token is invalid")
	}

	claims, ok := token.Claims.(*authTokenCustomClaims)

	if !ok {
		return nil, fmt.Errorf("ID token valid but couldn't parse claims")
	}

	return claims, nil
}

// validateRefreshToken uses the secret key to validate a refresh token
func validateRefreshToken(tokenString string, key string) (*refreshTokenCustomClaims, error) {
	claims := &refreshTokenCustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	// For now we'll just return the error and handle logging in service level
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, &apperror.Error{Message: apperror.InvalidRefreshToken, Type: apperror.Authorization}
	}

	claims, ok := token.Claims.(*refreshTokenCustomClaims)

	if !ok {
		return nil, &apperror.Error{Message: apperror.ErrClaimParse, Type: apperror.Authorization}
	}

	return claims, nil
}

// generateEmailVerficationOTP creates an otp
func generateEmailVerficationOTP(exp int64, from int, max int) (*model.OTPData, error) {
	currentTime := time.Now()
	rand.Seed(currentTime.UnixNano())

	result := rand.Intn(max-from) + from
	otpExp := currentTime.Add(time.Duration(exp) * time.Second)
	return &model.OTPData{
		OTP:       strconv.Itoa(result),
		ExpiresAt: otpExp,
	}, nil
}

func generateForgotPasswordOTP(exp int64, from int, max int) (*model.OTPData, error) {
	currentTime := time.Now()
	rand.Seed(currentTime.UnixNano())

	result := rand.Intn(max-from) + from
	otpExp := currentTime.Add(time.Duration(exp) * time.Second)
	return &model.OTPData{
		OTP:       strconv.Itoa(result),
		ExpiresAt: otpExp,
	}, nil
}

// validateEmailVerificationToken uses the secret key to validate a refresh token
func validateEmailVerificationToken(tokenString string, key string) (*jwt.StandardClaims, error) {
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	// For now we'll just return the error and handle logging in service level
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("refresh token is invalid")
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)

	if !ok {
		return nil, fmt.Errorf("refresh token valid but couldn't parse claims")
	}
	return claims, nil
}
