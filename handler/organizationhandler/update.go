package organizationhandler

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

func (h *HandlerImpl) Update() gin.HandlerFunc {
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

		organizationID := c.Param("organization_id")

		orgID, err := uuid.Parse(organizationID)
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

		var orgReq model.Organization

		if response, ok := handler.BindData2(c, &orgReq); !ok {
			return handler.HandlerResponse{
				Ctx:             reqCtx,
				ResponseWrapper: response,
				TrxKeys:         trxKeys,
				Ok:              false,
			}
		}

		// userTrx := h.UserService.BeginUserTrx()
		// reqCtx = context.WithValue(reqCtx, db.DBTrxUserKey, userTrx)
		// update updated_by value to current user
		getUser, err := h.UserService.Get(reqCtx, currUID)
		if err != nil {

			_, errResp, apperr := appresponse.PrepareErr(err, helper.TraceCurrentFuncArgs(map[string]string{
				"current_user_id": currUID.String(),
			}, orgReq))

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

		if getUser.Role != model.TSUPERADMIN && getUser.ID != currUID {
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

		err = h.OrganizationService.Update(reqCtx, orgID, &orgReq)
		if err != nil {

			_, errResp, apperr := appresponse.PrepareErr(err, helper.TraceCurrentFuncArgs(map[string]string{
				"current_user_id": currUID.String(),
			}, orgReq))

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
