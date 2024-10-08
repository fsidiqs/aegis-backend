package userhandler

import (
	"net/http"

	"github.com/fsidiqs/aegis-backend/handler"
	"github.com/fsidiqs/aegis-backend/helper"
	"github.com/fsidiqs/aegis-backend/model"
	"github.com/fsidiqs/aegis-backend/model/apperror"
	"github.com/fsidiqs/aegis-backend/model/appresponse"
	"github.com/gin-gonic/gin"
	"github.com/kr/pretty"
)

func (h *HandlerImpl) Register() gin.HandlerFunc {
	return handler.HandlerResolver(func(c *gin.Context) handler.HandlerResponse {
		var err error
		reqCtx := c.Request.Context()
		trxKeys := []string{}
		var ureq model.UserRegisterReq

		// publicUserVal, publicUserOk := c.Get("public_user")
		// if !publicUserOk {
		// 	errMsg := fmt.Sprintf("%s: failed get key-value: public_user", helper.TraceCurrentFunc())
		// 	// h.log.Error(errMsg)

		// 	return handler.HandlerResponse{
		// 		Ctx: reqCtx,
		// 		ResponseWrapper: appresponse.ResponseWrapper{
		// 			StatusCode: http.StatusInternalServerError,
		// 			Response: appresponse.ErrorResponse{
		// 				Message: appresponse.HdlMsgFailedExtractUserCtx, Type: appresponse.THandlerInternal,
		// 			},
		// 		}, TrxKeys: trxKeys, Ok: false,
		// 	}
		// }
		// publicUser := publicUserVal.(*model.ValidatedPublicToken)

		bindResp, ok := handler.BindData2(c, &ureq)
		if !ok {
			return handler.HandlerResponse{
				Ctx:             reqCtx,
				ResponseWrapper: bindResp, TrxKeys: trxKeys, Ok: false,
			}
		}

		u, err := ureq.CreateModel()
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

			return handler.HandlerResponse{
				Ctx: reqCtx,
				ResponseWrapper: appresponse.ResponseWrapper{
					StatusCode: sc,
					Response:   errResp,
				}, TrxKeys: trxKeys, Ok: false,
			}
		}
		pretty.Println("request", u)
		_, _, err = h.UserService.Register(reqCtx, *u)
		if err != nil {
			_, errResp, apperr := appresponse.PrepareErr(err, helper.TraceCurrentFuncArgs(ureq))
			sc := http.StatusInternalServerError

			if apperr.Type == apperror.BadRequest || apperr.Type == apperror.NotFound {
				// h.log.Info(internalMsg)
				sc = http.StatusBadRequest
				errResp = appresponse.HdlRespBadRequest()
			} else if apperr.Type == apperror.InternalConflict {
				sc = http.StatusBadRequest
				errResp = appresponse.ErrorResponse{Message: appresponse.HdlMsgUserExist, Type: appresponse.THdlUserExist}
			} else {
				// h.log.Error(internalMsg)
				sc = http.StatusInternalServerError
				errResp = appresponse.HdlRespInternalServerError()
			}

			return handler.HandlerResponse{
				Ctx: reqCtx,
				ResponseWrapper: appresponse.ResponseWrapper{
					StatusCode: sc,
					Response:   errResp,
				}, TrxKeys: trxKeys, Ok: false,
			}
		}

		// err = h.MailClient.SendEmailVerification(context.Background(), otp.OTP, userCreated.Email)
		err = h.MailClient.SendAccountCreatedMail(reqCtx, ureq.Password, u.Email)
		return handler.HandlerResponse{
			Ctx: reqCtx,
			ResponseWrapper: appresponse.ResponseWrapper{
				StatusCode: http.StatusOK,
				Response:   appresponse.SuccessResponse{Type: appresponse.THdlSuccess, Message: "user registered"},
			}, TrxKeys: trxKeys, Ok: true,
		}
	})

	// Bind incoming json to struct and check for validation errors

	// #region find referral code

	// #endregion find referral code
}
