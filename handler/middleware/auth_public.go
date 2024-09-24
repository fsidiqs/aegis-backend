package middleware

import (
	"net/http"
	"strings"

	"github.com/fsidiqs/aegis-backend/model/apperror"
	"github.com/fsidiqs/aegis-backend/model/appresponse"
	"github.com/fsidiqs/aegis-backend/service/tokenservice"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthPublicMiddleware struct {
	PublicTokenService tokenservice.IPublicTokenService
}

func (mid *AuthPublicMiddleware) AuthPublic() gin.HandlerFunc {
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

				c.JSON(err.Status(), appresponse.ErrorResponseMessageArr{Messages: apperror.InvalidArgumentsMap(invalidArgs), Type: appresponse.THdlBadRequest})
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
		validatedData, err := mid.PublicTokenService.ValidatePublicToken(authTokenHeader[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, appresponse.ErrorResponse{Message: appresponse.HdlMsgUnauthorized, Type: appresponse.THdlUnauthorized})
			c.Abort()

			return
		}

		c.Set("public_user", validatedData)
		c.Next()
		return
	}
}
