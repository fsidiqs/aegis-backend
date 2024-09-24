package handler

import (
	"net/http"

	"github.com/fsidiqs/aegis-backend/model/apperror"
	"github.com/fsidiqs/aegis-backend/model/appresponse"
	"github.com/gin-gonic/gin"
)

func ApperrorHandler(c *gin.Context, err *apperror.Error) {
	switch err.Type {
	case apperror.TFreeExceedMaximum:
		c.JSON(http.StatusBadRequest, appresponse.ErrorResponse{Message: appresponse.HdlMsgFreeUserMaximumEnroll, Type: appresponse.THdlBadRequest})
		return
	case apperror.BadRequest:
		c.JSON(http.StatusBadRequest, appresponse.ErrorResponse{Message: appresponse.HdlMsgBadRequest, Type: appresponse.THdlBadRequest})
		return
	default:
		c.JSON(http.StatusInternalServerError, appresponse.ErrorResponse{Message: appresponse.HdlMsgInternal, Type: appresponse.THandlerInternal})
	}
}
