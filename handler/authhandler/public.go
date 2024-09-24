package authhandler

import (
	"net/http"

	"github.com/fsidiqs/aegis-backend/handler"
	"github.com/fsidiqs/aegis-backend/helper"
	"github.com/fsidiqs/aegis-backend/model"
	"github.com/fsidiqs/aegis-backend/model/apperror"
	"github.com/fsidiqs/aegis-backend/model/appresponse"
	"github.com/gin-gonic/gin"
)

func (h *HandlerImpl) CreatePublicToken(c *gin.Context) {
	var req model.ReqPublicSession

	if ok := handler.BindData(c, &req); !ok {
		return
	}
	ctx := c.Request.Context()
	publicTok, err := h.PubTokenService.NewPairFromDeviceID(ctx, req.DeviceID, req.NotificationToken)
	if err != nil {

		_, errResp, apperr := appresponse.PrepareErr(err, helper.TraceCurrentFuncArgs(req))
		switch apperr.Type {
		case apperror.BadRequest:
			errResp = appresponse.HdlRespBadRequest()
		default:
			errResp = appresponse.HdlRespInternalServerError()
		}

		c.JSON(apperr.Status(), errResp)
		return
	}

	c.JSON(http.StatusOK, appresponse.SuccessResponse{Data: publicTok})
}
