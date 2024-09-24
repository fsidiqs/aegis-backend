package userhandler

import (
	"net/http"

	"github.com/fsidiqs/aegis-backend/handler"
	"github.com/fsidiqs/aegis-backend/helper"
	"github.com/fsidiqs/aegis-backend/model"
	"github.com/fsidiqs/aegis-backend/model/apperror"
	"github.com/fsidiqs/aegis-backend/model/appresponse"
	"github.com/gin-gonic/gin"
)

func (h *HandlerImpl) ForgotPasswordByEmail(c *gin.Context) {
	var ureq model.ReqUserResetPassword

	if ok := handler.BindData(c, &ureq); !ok {
		return
	}

	u, _ := ureq.ToModel()

	ctx := c.Request.Context()
	// perform UserService Register with transaction
	// trx := h.UserService.WithTrx()
	// perform register with created transaction
	err := h.UserService.ForgotPasswordUsingEmail(ctx, u.Email)
	if err != nil {

		_, errResp, apperr := appresponse.PrepareErr(err, helper.TraceCurrentFuncArgs(ureq))
		sc := http.StatusInternalServerError

		if apperr.Type == apperror.BadRequest {
			// h.log.Info(internalMsg)
			sc = http.StatusBadRequest
			errResp = appresponse.HdlRespBadRequest()
		} else if apperr.Type == apperror.TUserNotFound {
			sc = http.StatusNotFound
			errResp = appresponse.ErrorResponse{Message: appresponse.HdlMsgNoUser}
		} else {
			// h.log.Error(internalMsg)
			sc = http.StatusInternalServerError

			errResp = appresponse.HdlRespInternalServerError()
		}

		c.JSON(sc, errResp)
		return
	}

	c.Status(http.StatusOK)
}
