package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/fsidiqs/aegis-backend/model"
	"github.com/fsidiqs/aegis-backend/model/apperror"
	"github.com/fsidiqs/aegis-backend/model/appresponse"
	"github.com/fsidiqs/aegis-backend/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type authHeader struct {
	AuthToken string `header:"Authorization"`
}

type AuthMiddleware struct {
	TokenService service.ITokenService
	UserService  service.IUserService
	// Log          logservice.Logger
}

// AuthUser extracts a user from the Authorization header
// which is of the form "Bearer token"
// It sets the user to the context if the user exists
func (mid *AuthMiddleware) AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := authHeader{}

		// bind Authorization Header to h and check for validation errors
		if err := c.ShouldBindHeader(&h); err != nil {
			if errs, ok := err.(validator.ValidationErrors); ok {
				// we used this type in bind_data to extract desired fields from errs
				// you might consider extracting it
				var invalidArgs []apperror.InvalidArgument

				for _, err := range errs {
					invalidArgs = append(invalidArgs, apperror.InvalidArgument{
						Field: err.Field(),
						Value: err.Value().(string),
						Tag:   err.Tag(),
						Param: err.Param(),
					})
				}

				err := apperror.NewBadRequest("Invalid request parameters. See invalidArgs")

				c.JSON(err.Status(), gin.H{
					"error":       err,
					"invalidArgs": invalidArgs,
				})
				c.Abort()
				return
			}
		}

		authTokenHeader := strings.Split(h.AuthToken, "Bearer ")

		if len(authTokenHeader) < 2 {
			c.JSON(http.StatusUnauthorized, appresponse.ErrorResponse{Message: appresponse.HdlMsgNoAuthHeader, Type: appresponse.THdlNoAuthHeader})
			c.Abort()
			return
		}

		// 1. validate id token
		// 2. get refresh_token_id from extracted validated id token
		// 3. get session in redis session and perform check, if session exists then perform validation, if passes validation meanings client can continue, if no go to step 4
		// 4. if any error occurs while trying to fetch from redis or didnt pass validation then ignore and get session from db

		// validate ID token here
		validatedData, err := mid.TokenService.ValidateAuthToken(authTokenHeader[1])
		if err != nil {

			c.JSON(http.StatusUnauthorized, appresponse.ErrorResponse{Message: appresponse.HdlMsgUnauthorized, Type: appresponse.THdlUnauthorized})
			c.Abort()
			return
		}
		// ctx := c.Request.Context()
		// parse refresh token to uuid
		// reftokenid, err := uuid.Parse(validatedData.RefreshTokenID)
		_, err = uuid.Parse(validatedData.RefreshTokenID)
		if err != nil {
			// failed to parse refresh token id
			log.Printf("failed to parse refresh token id in middleware: %v\n", err)
			c.JSON(http.StatusInternalServerError, appresponse.ErrorResponse{Message: appresponse.HdlMsgInternal, Type: appresponse.THandlerInternal})

			c.Abort()
			return
		}
		// TODO CLEAN this was used before using redis or db fetch, it was only parsing from auth token
		// if user.AccountType == model.TUserUnverified {
		// 	c.JSON(http.StatusUnauthorized, appresponse.ErrorResponse{Message: appresponse.HdlMsgUserUnverified, Type: appresponse.THdlUnverified})
		// 	c.Abort()
		// 	return
		// }

		// if user.RecordFlag != model.Active {
		// 	c.JSON(http.StatusUnauthorized, appresponse.ErrorResponse{Message: appresponse.HdlMsgUnauthorized, Type: appresponse.THdlUnauthorized})
		// 	c.Abort()
		// 	return
		// }
		// redSess, err := mid.UserService.RedisGetSession(ctx, validatedData.User.ID, reftokenid)
		// // if no error then return the current redis session
		// // if there is error, just log and ignore, and continue find session in db
		// if err != nil {
		// mid.Log.Info("failed get session from user redis, proceeding to fetch user session from db: %v", h.AuthToken)
		// }

		// if redSess != nil && redSess.UserRecordFlag == model.RecActive && redSess.AccountType != model.TUserUnverified {
		// 	c.Set("user", validatedData.User)
		// 	c.Next()
		// 	return
		// }

		// dbSess, err := mid.UserService.FindUserSessionDBByUserID(ctx, validatedData.User.ID)
		// if err != nil {
		// 	// failed to parse refresh token id
		// 	// mid.Log.Info("error fetching from db:userid:%v: err: %v\n", validatedData.User.ID, err)
		// 	c.JSON(http.StatusUnauthorized, appresponse.ErrorResponse{Message: appresponse.HdlMsgUnauthorized, Type: appresponse.THdlUnauthorized})

		// 	c.Abort()
		// 	return
		// }

		// if dbSess == nil {
		// 	// mid.Log.Info("no available session from db :userid:%v: err: %v\n", validatedData.User.ID, err)

		// 	c.JSON(http.StatusUnauthorized, appresponse.ErrorResponse{Message: appresponse.HdlMsgUnauthorized, Type: appresponse.THdlUnauthorized})
		// 	c.Abort()
		// 	return
		// }
		// if dbSess.AccountType == model.TUserUnverified {
		// 	c.JSON(http.StatusUnauthorized, appresponse.ErrorResponse{Message: appresponse.HdlMsgUserUnverified, Type: appresponse.THdlUnverified})
		// 	c.Abort()
		// 	return
		// }
		// if dbSess.AuthToken != authTokenHeader[1] {
		// 	c.JSON(http.StatusUnauthorized, appresponse.ErrorResponse{Message: appresponse.HdlMsgUnauthorized, Type: appresponse.THdlUnauthorized})
		// 	c.Abort()
		// }

		dbsessExtract, err := mid.TokenService.ValidateAuthToken(authTokenHeader[1])
		if err != nil {
			// failed to parse refresh token id
			// mid.Log.Info("error validating and extracting from session:userid:%v: err: %v\n", validatedData.User.ID, err)
			c.JSON(http.StatusInternalServerError, appresponse.ErrorResponse{Message: appresponse.HdlMsgInternal, Type: appresponse.THandlerInternal})

			c.Abort()
			return
		}

		if dbsessExtract.User.RecordFlag != model.RecActive {
			// mid.Log.Info("user session is not active: %v\n", dbsessExtract.User.ID)
			c.JSON(http.StatusUnauthorized, appresponse.ErrorResponse{Message: appresponse.HdlMsgUserNotActive, Type: appresponse.THdlUnauthorized})
			c.Abort()
			return
		}

		// if dbSess.Status == model.TUserSessionLogout {
		// 	c.JSON(http.StatusUnauthorized, appresponse.ErrorResponse{Message: appresponse.HdlMsgUnauthorized, Type: appresponse.THdlUnauthorized})
		// 	c.Abort()
		// 	return
		// }
		c.Set("user", validatedData.User)
		c.Next()
		return
	}
}

func (mid *AuthMiddleware) AuthInternal() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := authHeader{}

		// bind Authorization Header to h and check for validation errors
		if err := c.ShouldBindHeader(&h); err != nil {
			if errs, ok := err.(validator.ValidationErrors); ok {
				// we used this type in bind_data to extract desired fields from errs
				// you might consider extracting it
				var invalidArgs []apperror.InvalidArgument

				for _, err := range errs {
					invalidArgs = append(invalidArgs, apperror.InvalidArgument{
						Field: err.Field(),
						Value: err.Value().(string),
						Tag:   err.Tag(),
						Param: err.Param(),
					})
				}

				err := apperror.NewBadRequest("Invalid request parameters. See invalidArgs")

				c.JSON(err.Status(), gin.H{
					"error":       err,
					"invalidArgs": invalidArgs,
				})
				c.Abort()
				return
			}
		}

		authTokenHeader := strings.Split(h.AuthToken, "Bearer ")

		if len(authTokenHeader) < 2 {
			c.JSON(http.StatusUnauthorized, appresponse.ErrorResponse{Message: appresponse.HdlMsgNoAuthHeader, Type: appresponse.THdlNoAuthHeader})
			c.Abort()
			return
		}

		// 1. validate id token
		// 2. get refresh_token_id from extracted validated id token
		// 3. get session in redis session and perform check, if session exists then perform validation, if passes validation meanings client can continue, if no go to step 4
		// 4. if any error occurs while trying to fetch from redis or didnt pass validation then ignore and get session from db

		// validate ID token here
		validatedData, err := mid.TokenService.ValidateAuthToken(authTokenHeader[1])
		if err != nil {

			c.JSON(http.StatusUnauthorized, appresponse.ErrorResponse{Message: appresponse.HdlMsgUnauthorized, Type: appresponse.THdlUnauthorized})
			c.Abort()
			return
		}
		ctx := c.Request.Context()
		// parse refresh token to uuid

		// dbSess, err := mid.UserService.FindUserSessionDBByUserID(ctx, validatedData.User.ID)
		// if err != nil {
		// 	// failed to parse refresh token id
		// 	// mid.Log.Info("error fetching from db:userid:%v: err: %v\n", validatedData.User.ID, err)
		// 	c.JSON(http.StatusUnauthorized, appresponse.ErrorResponse{Message: appresponse.HdlMsgUnauthorized, Type: appresponse.THdlUnauthorized})

		// 	c.Abort()
		// 	return
		// }

		// if dbSess == nil {
		// 	// mid.Log.Info("no available session from db :userid:%v: err: %v\n", validatedData.User.ID, err)

		// 	c.JSON(http.StatusUnauthorized, appresponse.ErrorResponse{Message: appresponse.HdlMsgUnauthorized, Type: appresponse.THdlUnauthorized})
		// 	c.Abort()
		// 	return
		// }
		// if dbSess.AccountType == model.TUserUnverified {
		// 	c.JSON(http.StatusUnauthorized, appresponse.ErrorResponse{Message: appresponse.HdlMsgUserUnverified, Type: appresponse.THdlUnverified})
		// 	c.Abort()
		// 	return
		// }
		// if dbSess.AuthToken != authTokenHeader[1] {
		// 	c.JSON(http.StatusUnauthorized, appresponse.ErrorResponse{Message: appresponse.HdlMsgUnauthorized, Type: appresponse.THdlUnauthorized})
		// 	c.Abort()
		// }

		dbsessExtract, err := mid.TokenService.ValidateAuthToken(authTokenHeader[1])
		if err != nil {
			// failed to parse refresh token id
			// mid.Log.Info("error validating and extracting from session:userid:%v: err: %v\n", validatedData.User.ID, err)
			c.JSON(http.StatusInternalServerError, appresponse.ErrorResponse{Message: appresponse.HdlMsgInternal, Type: appresponse.THandlerInternal})

			c.Abort()
			return
		}

		if dbsessExtract.User.RecordFlag != model.RecActive {
			// mid.Log.Info("user session is not active: %v\n", dbsessExtract.User.ID)
			c.JSON(http.StatusUnauthorized, appresponse.ErrorResponse{Message: appresponse.HdlMsgUserNotActive, Type: appresponse.THdlUnauthorized})
			c.Abort()
			return
		}

		// if dbSess.Status == model.TUserSessionLogout {
		// 	c.JSON(http.StatusUnauthorized, appresponse.ErrorResponse{Message: appresponse.HdlMsgUnauthorized, Type: appresponse.THdlUnauthorized})
		// 	c.Abort()
		// 	return
		// }

		// validate if it is internal user
		userFetch, err := mid.UserService.Get(ctx, validatedData.User.ID)
		if err != nil {
			// mid.Log.Info("error fetching user:%v\n", err)
			c.JSON(http.StatusInternalServerError, appresponse.ErrorResponse{Message: appresponse.HdlMsgInternal, Type: appresponse.THandlerInternal})
			c.Abort()
			return
		}

		if userFetch.Email != "fajar@mindtera.com" {
			c.JSON(http.StatusBadRequest, appresponse.ErrorResponse{Message: appresponse.HdlMsgUnauthorized, Type: appresponse.THdlUnauthorized})
			c.Abort()
			return
		}

		c.Set("user", validatedData.User)
		c.Next()
		return
	}
}
