package organizationhandler

import (
	"net/http"

	"github.com/fsidiqs/aegis-backend/handler"
	"github.com/fsidiqs/aegis-backend/helper"
	"github.com/fsidiqs/aegis-backend/model"
	"github.com/fsidiqs/aegis-backend/model/apperror"
	"github.com/fsidiqs/aegis-backend/model/appresponse"
	"github.com/gin-gonic/gin"
)

func (h *HandlerImpl) SubmitOTPForgotPassword(c *gin.Context) {
	var req model.ReqSubmitOTPForgotPassword

	if ok := handler.BindData(c, &req); !ok {
		return
	}

	ctx := c.Request.Context()
	// perform UserService Register with transaction

	// perform update password service
	err := h.UserService.SubmitOTPForgotPassword(ctx, req.Email, req.OTP)
	if err != nil {

		_, errResp, apperr := appresponse.PrepareErr(err, helper.TraceCurrentFuncArgs(req))
		sc := http.StatusInternalServerError
		if apperr.Type == apperror.BadRequest || apperr.Type == apperror.NotFound {
			sc = http.StatusBadRequest
			errResp = appresponse.HdlRespBadRequest()
		} else if apperr.Type == apperror.EmailNotFound {
			sc = http.StatusBadRequest
			errResp = appresponse.ErrorResponse{Message: appresponse.MsgEmailNotFound, Type: appresponse.TEmailNotFound}
		} else if apperr.Type == apperror.OTPNotFound {
			sc = http.StatusBadRequest
			errResp = appresponse.ErrorResponse{Message: appresponse.MsgOTPNotFound, Type: appresponse.TOTPNotFound}
		} else if apperr.Type == apperror.OTPExpired {
			sc = http.StatusBadRequest
			errResp = appresponse.ErrorResponse{Message: appresponse.HdlMsgOTPExpire, Type: appresponse.TOTPExpired}
		} else {
			sc = http.StatusInternalServerError
			errResp = appresponse.HdlRespInternalServerError()
		}

		c.JSON(sc, errResp)
		return

	}

	c.Status(http.StatusOK)
}
