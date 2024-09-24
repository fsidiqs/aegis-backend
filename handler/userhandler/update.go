package userhandler

import (
	"net/http"

	"github.com/fsidiqs/aegis-backend/db"
	"github.com/fsidiqs/aegis-backend/handler"
	"github.com/fsidiqs/aegis-backend/helper"
	"github.com/fsidiqs/aegis-backend/model"
	"github.com/fsidiqs/aegis-backend/model/apperror"
	"github.com/fsidiqs/aegis-backend/model/appresponse"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kr/pretty"
)

func (h *HandlerImpl) UpdateMyDetails() gin.HandlerFunc {
	return handler.HandlerResolver(func(c *gin.Context) handler.HandlerResponse {
		reqCtx := c.Request.Context()
		trxKeys := []string{}

		// get user context extracted by auth_user middleware
		userVal, ok := c.Get("user")

		if !ok {
			return handler.HandlerResponse{
				Ctx: reqCtx,
				ResponseWrapper: appresponse.ResponseWrapper{
					StatusCode: http.StatusInternalServerError,
					Response:   appresponse.ErrorResponse{Message: "failed to extract user"},
				},
				TrxKeys: trxKeys,
				Ok:      false,
			}
		}
		user := userVal.(*model.User)
		var req model.UserUpdateReq

		if response, ok := handler.BindData2(c, &req); !ok {
			return handler.HandlerResponse{
				Ctx:             reqCtx,
				ResponseWrapper: response,
				TrxKeys:         trxKeys,
				Ok:              false,
			}
		}

		// perform UserService Update
		uUpdate, err := req.ToUserUpdateModel(user.ID.String())
		if err != nil {

			_, errResp, apperr := appresponse.PrepareErr(err, helper.TraceCurrentFuncArgs(map[string]string{
				"user_id": user.ID.String(),
			}, req))

			sc := http.StatusInternalServerError
			if apperr.Type == apperror.BadRequest {
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

		if err := h.UserService.UpdateUser(reqCtx, user.ID, uUpdate); err != nil {
			_, errResp, apperr := appresponse.PrepareErr(err, helper.TraceCurrentFuncArgs(map[string]string{
				"user_id": user.ID.String(),
			}, req))

			sc := http.StatusInternalServerError
			if apperr.Type == apperror.BadRequest {
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
		trxKeys = append(trxKeys, db.DBTrxUserKey)
		// user_activity_log

		trxKeys = append(trxKeys, db.DBTrxUserActivityLogKey)

		return handler.HandlerResponse{
			Ctx: reqCtx,
			ResponseWrapper: appresponse.ResponseWrapper{
				StatusCode: http.StatusOK,
				Response:   appresponse.SuccessResponse{Type: appresponse.THdlSuccess, Message: "user updated successfully"},
			}, TrxKeys: trxKeys, Ok: true,
		}
	})
}

func (h *HandlerImpl) UpdateDetails() gin.HandlerFunc {
	return handler.HandlerResolver(func(c *gin.Context) handler.HandlerResponse {
		reqCtx := c.Request.Context()
		trxKeys := []string{}

		currUser, ok := c.Get("user")

		if !ok {
			return handler.HandlerResponse{
				Ctx: reqCtx,
				ResponseWrapper: appresponse.ResponseWrapper{
					StatusCode: http.StatusInternalServerError,
					Response:   appresponse.ErrorResponse{Message: "failed to extract user"},
				},
				TrxKeys: trxKeys,
				Ok:      false,
			}
		}
		currUID := currUser.(*model.User).ID

		user := c.Param("user_id")
		pretty.Println("user_id: ", currUser)
		pretty.Println()
		uid, err := uuid.Parse(user)
		if err != nil {
			return handler.HandlerResponse{
				Ctx: reqCtx,
				ResponseWrapper: appresponse.ResponseWrapper{
					StatusCode: http.StatusBadRequest,
					Response: appresponse.ErrorResponse{
						Type: appresponse.THdlBadRequest, Message: appresponse.HdlMsgBadRequest,
					},
				}, TrxKeys: trxKeys, Ok: false,
			}
		}

		var req model.UserUpdateReq

		if response, ok := handler.BindData2(c, &req); !ok {
			return handler.HandlerResponse{
				Ctx:             reqCtx,
				ResponseWrapper: response,
				TrxKeys:         trxKeys,
				Ok:              false,
			}
		}

		// perform UserService Update
		uUpdate, err := req.ToUserUpdateModel(uid.String())
		if err != nil {

			_, errResp, apperr := appresponse.PrepareErr(err, helper.TraceCurrentFuncArgs(map[string]string{
				"current_user_id": currUID.String(),
				"user_id":         uid.String(),
			}, req))

			sc := http.StatusInternalServerError
			if apperr.Type == apperror.BadRequest {
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

		// userTrx := h.UserService.BeginUserTrx()
		// reqCtx = context.WithValue(reqCtx, db.DBTrxUserKey, userTrx)
		// update updated_by value to current user
		getUser, err := h.UserService.Get(reqCtx, currUID)
		if err != nil {

			_, errResp, apperr := appresponse.PrepareErr(err, helper.TraceCurrentFuncArgs(map[string]string{
				"current_user_id": currUID.String(),
				"user_id":         uid.String(),
			}, req))

			sc := http.StatusInternalServerError
			if apperr.Type == apperror.BadRequest {
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

		if getUser.Role != model.TSUPERADMIN && getUser.ID != uid {
			pretty.Println("getUser.Role: ", getUser)
			pretty.Println("getUser.ID: ", getUser.ID)
			return handler.HandlerResponse{
				Ctx: reqCtx,
				ResponseWrapper: appresponse.ResponseWrapper{
					StatusCode: http.StatusForbidden,
					Response:   appresponse.ErrorResponse{Message: "forbidden access"},
				}, TrxKeys: trxKeys, Ok: false,
			}
		}

		if err := h.UserService.UpdateUser(reqCtx, uid, uUpdate); err != nil {
			_, errResp, apperr := appresponse.PrepareErr(err, helper.TraceCurrentFuncArgs(map[string]string{
				"current_user_id": currUID.String(),
				"user_id":         uid.String(),
			}, req))

			sc := http.StatusInternalServerError
			if apperr.Type == apperror.BadRequest {
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
		trxKeys = append(trxKeys, db.DBTrxUserKey)

		// user_activity_log
		// actLogTrx := h.ActivityLogService.BeginTrx()
		// reqCtx = context.WithValue(reqCtx, db.DBTrxUserActivityLogKey, actLogTrx)

		// h.ActivityLogService.Store(reqCtx, model.UserActivityLog{
		// UserId:            currUID,
		// ActivityDate:      time.Now(),
		// ActivityType:      "User",
		// RelatedActivityID: &uid,
		// DefaultColumns:    model.NewLogDefaultCol(helper.TraceCurrentFunc()),
		// })
		trxKeys = append(trxKeys, db.DBTrxUserActivityLogKey)

		return handler.HandlerResponse{
			Ctx: reqCtx,
			ResponseWrapper: appresponse.ResponseWrapper{
				StatusCode: http.StatusOK,
				Response:   appresponse.SuccessResponse{Type: appresponse.THdlSuccess, Message: "user updated successfully"},
			}, TrxKeys: trxKeys, Ok: true,
		}
	})
}
