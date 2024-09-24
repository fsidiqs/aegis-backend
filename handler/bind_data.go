package handler

import (
	"net/http"

	"github.com/fsidiqs/aegis-backend/model/apperror"
	"github.com/fsidiqs/aegis-backend/model/appresponse"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type errorResp struct {
	Err         *apperror.Error            `json:"error"`
	InvalidArgs []apperror.InvalidArgument `json:"invalidArgs"`
}

// BindData is helper function, returns false if data is not found

func BindData(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBind(req); err != nil {
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

			err := apperror.NewBadRequest(string(appresponse.HdlMsgBadRequest))

			c.JSON(err.Status(), gin.H{
				"type":         err.Type,
				"message":      err.Message,
				"invalid_args": invalidArgs,
			})
			return false
		}
	}
	return true
}

func BindData2(c *gin.Context, req interface{}) (appresponse.ResponseWrapper, bool) {
	if err := c.ShouldBind(req); err != nil {
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

			return appresponse.ResponseWrapper{
				StatusCode: http.StatusBadRequest,
				Response: appresponse.ErrRespWithInvalidArgs{
					Type:        appresponse.THdlBadRequest,
					Message:     appresponse.HdlMsgBadRequest,
					InvalidArgs: invalidArgs,
				},
			}, false
		}
	}
	return appresponse.ResponseWrapper{}, true
}
