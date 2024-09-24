package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/fsidiqs/aegis-backend/model"
	"github.com/fsidiqs/aegis-backend/model/apperror"
	"github.com/fsidiqs/aegis-backend/model/appresponse"
	"github.com/fsidiqs/aegis-backend/service"
	"github.com/fsidiqs/aegis-backend/service/tokenservice"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// this middleware has non aborting (will allow next handler being processed)while validating jwt token,
// purpose: try validating jwt token using public_auth when validating using auth_user fails

type AuthUserOrPublicMiddleware struct {
	PublicTokenService tokenservice.IPublicTokenService
	TokenService       service.ITokenService
	UserService        service.IUserService
}

func (mid *AuthUserOrPublicMiddleware) AuthUserOrAuthPublic() gin.HandlerFunc {
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
		// TODO this checking is costly, probably should another create Authorization header format e.g: "Bearer smth <token>"
		if err != nil {
			validatedData, err := mid.performPublicTokenValidation(authTokenHeader[1])
			if err != nil {
				c.JSON(http.StatusUnauthorized, appresponse.ErrorResponse{Message: appresponse.HdlMsgUnauthorized, Type: appresponse.THdlUnauthorized})

				c.Abort()
				return
			}
			c.Set("public_user", validatedData)
			c.Next()
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

		// redSess, err := mid.UserService.RedisGetSession(ctx, validatedData.User.ID, reftokenid)
		// // if no error then return the current redis session
		// // if there is error, just log and ignore to continue find session in db
		// if err == nil {
		// 	if redSess.UserRecordFlag == model.RecActive && redSess.AccountType != model.TUserUnverified {
		// 		c.Set("user", validatedData.User)
		// 		c.Next()
		// 		return
		// 	}
		// } else {
		// 	log.Printf("error fetching redis:userid:%v err:%v", validatedData.User.ID, err)
		// }

		// dbSess, err := mid.UserService.FindUserSessionDBByUserID(ctx, validatedData.User.ID)
		// if err != nil {
		// 	// failed to parse refresh token id
		// 	log.Printf("error fetching from db:userid:%v: err: %v\n", validatedData.User.ID, err)
		// 	c.JSON(http.StatusUnauthorized, appresponse.ErrorResponse{Message: appresponse.HdlMsgUnauthorized, Type: appresponse.THdlUnauthorized})

		// 	c.Abort()
		// 	return
		// }

		// if dbSess == nil {
		// 	log.Printf("no available session from db :userid:%v: err: %v\n", validatedData.User.ID, err)

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
		// dbsessExtract, err := mid.TokenService.ValidateAuthToken(dbSess.AuthToken)
		// if err != nil {
		// 	// failed to parse refresh token id
		// 	log.Printf("error validating and extracting from session:userid:%v: err: %v\n", validatedData.User.ID, err)
		// 	c.JSON(http.StatusInternalServerError, appresponse.ErrorResponse{Message: appresponse.HdlMsgInternal, Type: appresponse.THandlerInternal})

		// 	c.Abort()
		// 	return
		// }

		// if dbsessExtract.User.RecordFlag != model.RecActive {
		// 	c.JSON(http.StatusUnauthorized, appresponse.ErrorResponse{Message: appresponse.HdlMsgUserNotActive, Type: appresponse.THdlUnauthorized})
		// 	c.Abort()
		// 	return
		// }

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

func (mid *AuthUserOrPublicMiddleware) performPublicTokenValidation(authToken string) (*model.ValidatedPublicToken, error) {
	validatedData, err := mid.PublicTokenService.ValidatePublicToken(authToken)
	if err != nil {
		return nil, &apperror.Error{Type: apperror.Authorization}
	}

	return validatedData, nil
}
