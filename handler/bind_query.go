package handler

import (
	"github.com/fsidiqs/aegis-backend/model/apperror"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// BindQuery is helper function, returns false if data is not found

func BindQuery(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindQuery(req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			// could probably extract this, it is also in middleware_auth_user
			var invalidArgs []apperror.InvalidArgument

			for _, err := range errs {
				invalidArgs = append(invalidArgs, apperror.InvalidArgument{
					Field: err.Field(),
					Value: err.Value().(string),
					Tag:   err.Tag(),
					Param: err.Param(),
				})
			}

			err := apperror.NewBadRequest("Invalid query params")

			c.JSON(err.Status(), errorResp{
				Err:         err,
				InvalidArgs: invalidArgs,
			})
			return false
		}
	}
	return true
}
