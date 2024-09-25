package organizationhandler

import (
	"log"
	"net/http"

	"github.com/fsidiqs/aegis-backend/handler"
	"github.com/fsidiqs/aegis-backend/helper"
	"github.com/fsidiqs/aegis-backend/model"
	"github.com/fsidiqs/aegis-backend/model/apperror"
	"github.com/fsidiqs/aegis-backend/model/appresponse"
	"github.com/gin-gonic/gin"
)

func (h *HandlerImpl) UpdateForgottenPasswordByOTP(c *gin.Context) {
	var ureq model.ReqUpdateForgottenPasswordByOTP

	if ok := handler.BindData(c, &ureq); !ok {
		return
	}
	u, err := ureq.ToUserPasswordUpdateModel(ureq.Email)
	if err != nil {
		log.Printf("Failed to create user object:email:%v err:%v\n", ureq.Email, err.Error())
		_, errResp, apperr := appresponse.PrepareErr(err, helper.TraceCurrentFuncArgs(ureq))
		sc := http.StatusInternalServerError

		if apperr.Type == apperror.BadRequest || apperr.Type == apperror.NotFound {
			// h.log.Info(internalMsg)
			sc = http.StatusBadRequest
			errResp = appresponse.HdlRespBadRequest()
		} else {
			// h.log.Error(internalMsg)
			sc = http.StatusInternalServerError
			errResp = appresponse.HdlRespInternalServerError()
		}

		c.JSON(sc, errResp)
		return

	}

	ctx := c.Request.Context()
	// perform UserService Register with transaction
	// trx := h.UserService.WithTrx()

	// perform update password service
	uFetch, err := h.UserService.UpdatePasswordUsingEmailAndOTP(ctx, ureq.Email, ureq.OTP, *u)
	if err != nil {
		_, errResp, apperr := appresponse.PrepareErr(err, helper.TraceCurrentFuncArgs(ureq))
		sc := http.StatusInternalServerError

		if apperr.Type == apperror.BadRequest || apperr.Type == apperror.NotFound {
			// h.log.Info(internalMsg)
			sc = http.StatusBadRequest
			errResp = appresponse.HdlRespBadRequest()
		} else {
			// h.log.Error(internalMsg)
			sc = http.StatusInternalServerError
			errResp = appresponse.HdlRespInternalServerError()
		}

		c.JSON(sc, errResp)
		return
	}

	tokens, err := h.TokenService.NewPairFromUser(ctx, uFetch)
	if err != nil {
		_, errResp, apperr := appresponse.PrepareErr(err, helper.TraceCurrentFuncArgs(ureq))
		sc := http.StatusInternalServerError

		if apperr.Type == apperror.BadRequest || apperr.Type == apperror.NotFound {
			// h.log.Info(internalMsg)
			sc = http.StatusBadRequest
			errResp = appresponse.HdlRespBadRequest()
		} else {
			// h.log.Error(internalMsg)
			sc = http.StatusInternalServerError
			errResp = appresponse.HdlRespInternalServerError()
		}

		c.JSON(sc, errResp)
		return
	}

	// BEGIN: creating session
	// uSess, err := h.UserService.FindUserSessionDBByUserID(ctx, uFetch.ID)
	// if err != nil {

	// 	_, errResp, apperr := appresponse.PrepareErr(err, helper.TraceCurrentFuncArgs(ureq))
	// 	sc := http.StatusInternalServerError

	// 	if apperr.Type == apperror.BadRequest || apperr.Type == apperror.NotFound {
	// 		// h.log.Info(internalMsg)
	// 		sc = http.StatusBadRequest
	// 		errResp = appresponse.HdlRespBadRequest()
	// 	} else {
	// 		// h.log.Error(internalMsg)
	// 		sc = http.StatusInternalServerError
	// 		errResp = appresponse.HdlRespInternalServerError()
	// 	}

	// 	c.JSON(sc, errResp)
	// 	return
	// }

	// if uSess != nil {
	// 	// now := time.Now()
	// 	// uSessUpdate := model.UserSessionUpdate{
	// 	// 	AuthToken:      tokens.AuthToken.SS,
	// 	// 	RefreshTokenID: tokens.RefreshToken.ID.String(),
	// 	// 	ExpiredAt:      now.Add(tokens.ExpiresIn),
	// 	// 	Status:         model.TUserSessionActive,
	// 	// 	DefaultColumns: model.DefaultColumns{
	// 	// 		UpdatedBy: "forgot-password-handler",
	// 	// 		UpdatedAt: &now,
	// 	// 	},
	// 	// }
	// 	// if err != nil {

	// 	// 	_, errResp, apperr := appresponse.PrepareErr(err, helper.TraceCurrentFuncArgs(ureq))
	// 	// 	sc := http.StatusInternalServerError

	// 	// 	if apperr.Type == apperror.BadRequest || apperr.Type == apperror.NotFound {
	// 	// 		// h.log.Info(internalMsg)
	// 	// 		sc = http.StatusBadRequest
	// 	// 		errResp = appresponse.HdlRespBadRequest()
	// 	// 	} else {
	// 	// 		// h.log.Error(internalMsg)
	// 	// 		sc = http.StatusInternalServerError
	// 	// 		errResp = appresponse.HdlRespInternalServerError()
	// 	// 	}

	// 	// 	c.JSON(sc, errResp)
	// 	// 	return
	// 	// }
	// } else {
	// 	// store session to db and redis

	// 	// uSess := model.NewUserSession(tokens,
	// 	// 	model.ExtraData{
	// 	// 		Status:         model.TUserSessionActive,
	// 	// 		AccountType:    uFetch.AccountType,
	// 	// 		UserRecordFlag: uFetch.RecordFlag,
	// 	// 	})

	// 	// if err != nil {

	// 	// 	_, errResp, apperr := appresponse.PrepareErr(err, helper.TraceCurrentFuncArgs(ureq))
	// 	// 	sc := http.StatusInternalServerError

	// 	// 	if apperr.Type == apperror.BadRequest || apperr.Type == apperror.NotFound {
	// 	// 		// h.log.Info(internalMsg)
	// 	// 		sc = http.StatusBadRequest
	// 	// 		errResp = appresponse.HdlRespBadRequest()
	// 	// 	} else {
	// 	// 		// h.log.Error(internalMsg)
	// 	// 		sc = http.StatusInternalServerError
	// 	// 		errResp = appresponse.HdlRespInternalServerError()
	// 	// 	}

	// 	// 	c.JSON(sc, errResp)
	// 	// 	return
	// 	// }
	// }

	// lastLoginMethod := &model.UserUpdate{
	// 	LastLoginMethod: string(model.TLastLoginManual),
	// 	DefaultColumns: model.DefaultColumns{
	// 		UpdatedBy: helper.TraceCurrentFunc(),
	// 	},
	// }
	// if err != nil {

	// 	_, errResp, apperr := appresponse.PrepareErr(err, helper.TraceCurrentFuncArgs(ureq))
	// 	sc := http.StatusInternalServerError

	// 	if apperr.Type == apperror.BadRequest || apperr.Type == apperror.NotFound {
	// 		// h.log.Info(internalMsg)
	// 		sc = http.StatusBadRequest
	// 		errResp = appresponse.HdlRespBadRequest()
	// 	} else {
	// 		// h.log.Error(internalMsg)
	// 		sc = http.StatusInternalServerError
	// 		errResp = appresponse.HdlRespInternalServerError()
	// 	}

	// 	c.JSON(sc, errResp)
	// 	return
	// }

	// user_activity_log

	// err = h.ActivityLogService.Store(ctx, model.UserActivityLog{
	// 	UserId:            uFetch.ID,
	// 	ActivityDate:      time.Now(),
	// 	ActivityType:      "users",
	// 	RelatedActivityID: nil,
	// 	DefaultColumns:    model.NewLogDefaultCol(helper.TraceCurrentFunc()),
	// })
	// if err != nil {

	// 	_, errResp, apperr := appresponse.PrepareErr(err, helper.TraceCurrentFuncArgs(ureq))
	// 	sc := http.StatusInternalServerError

	// 	if apperr.Type == apperror.BadRequest || apperr.Type == apperror.NotFound {
	// 		// h.log.Info(internalMsg)
	// 		sc = http.StatusBadRequest
	// 		errResp = appresponse.HdlRespBadRequest()
	// 	} else {
	// 		// h.log.Error(internalMsg)
	// 		sc = http.StatusInternalServerError
	// 		errResp = appresponse.HdlRespInternalServerError()
	// 	}

	// 	c.JSON(sc, errResp)
	// 	return
	// }

	c.JSON(http.StatusOK, appresponse.SuccessResponse{Data: tokens})
}
