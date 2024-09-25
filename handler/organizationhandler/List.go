package organizationhandler

import (
	"net/http"

	"github.com/fsidiqs/aegis-backend/handler"
	"github.com/fsidiqs/aegis-backend/helper"
	"github.com/fsidiqs/aegis-backend/model"
	"github.com/fsidiqs/aegis-backend/model/apperror"
	"github.com/fsidiqs/aegis-backend/model/appresponse"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RegisterUser
// @Summary Handling request from the client to perform user registration
// @Tags userRegister
// @Accept json
// @Produce json
// @Param user body model.UserRegisterReq true "contains user data, the combination of email and phone must be unique"
// @Success 200 "ok"
// @Failure 400 {object} apperror.ErrorResponse "bad request"
// @Router /consumer/v1/user/{user_id} [post]
func (h *HandlerImpl) MyDetails() gin.HandlerFunc {
	return handler.HandlerResolver(func(c *gin.Context) handler.HandlerResponse {
		userKey, ok := c.Get("user")
		if !ok {
			// errMsg := fmt.Sprintf("%s: failed get key-value: user", helper.TraceCurrentFunc())
			return handler.HandlerResponse{
				ResponseWrapper: appresponse.ResponseWrapper{
					StatusCode: http.StatusInternalServerError,
					Response:   appresponse.ErrorResponse{Type: appresponse.THdlBadRequest, Message: "failed to get user key-value"},
				},
				Ok: false,
			}
		}

		user := userKey.(*model.User)
		reqCtx := c.Request.Context()

		uFetched, err := h.UserService.Get(reqCtx, user.ID)
		if err != nil {
			_, errResp, apperr := appresponse.PrepareErr(err, helper.TraceCurrentFuncArgs(map[string]string{
				"user_id": user.ID.String(),
			}))
			sc := http.StatusInternalServerError

			if apperr.Type == apperror.BadRequest || apperr.Type == apperror.NotFound {
				sc = http.StatusBadRequest
				errResp = appresponse.HdlRespBadRequest()
			} else {
				sc = http.StatusInternalServerError
				errResp = appresponse.HdlRespInternalServerError()
			}

			return handler.HandlerResponse{
				ResponseWrapper: appresponse.ResponseWrapper{
					StatusCode: sc,
					Response:   errResp,
				}, Ok: false,
			}
		}

		return handler.HandlerResponse{
			ResponseWrapper: appresponse.ResponseWrapper{
				StatusCode: http.StatusOK,
				Response:   appresponse.SuccessResponse{Data: uFetched, Type: appresponse.THdlSuccess, Message: appresponse.HdlMsgSuccess},
			}, Ok: true,
		}
	})
}

func (h *HandlerImpl) Get() gin.HandlerFunc {
	return handler.HandlerResolver(func(c *gin.Context) handler.HandlerResponse {
		OrgIDParam := c.Param("organization_id")

		orgID, err := uuid.Parse(OrgIDParam)
		if err != nil {
			// errMsg := fmt.Sprintf("%s: failed parse user-uuid :%v", helper.TraceCurrentFunc(), userIDParam)
			return handler.HandlerResponse{
				ResponseWrapper: appresponse.ResponseWrapper{
					StatusCode: http.StatusInternalServerError,
					Response:   appresponse.ErrorResponse{Type: appresponse.THdlBadRequest, Message: "failed parse user-uuid"},
				},
				Ok: false,
			}
		}

		ctx := c.Request.Context()
		orgGet, err := h.OrganizationService.Get(ctx, orgID)
		if err != nil {
			_, errResp, apperr := appresponse.PrepareErr(err, helper.TraceCurrentFuncArgs(map[string]string{}))
			sc := http.StatusInternalServerError

			if apperr.Type == apperror.BadRequest || apperr.Type == apperror.NotFound {
				sc = http.StatusBadRequest
				errResp = appresponse.HdlRespBadRequest()
			} else {
				sc = http.StatusInternalServerError
				errResp = appresponse.HdlRespInternalServerError()
			}

			return handler.HandlerResponse{
				ResponseWrapper: appresponse.ResponseWrapper{
					StatusCode: sc,
					Response:   errResp,
				}, Ok: false,
			}
		}

		return handler.HandlerResponse{
			ResponseWrapper: appresponse.ResponseWrapper{
				StatusCode: http.StatusOK,
				Response:   appresponse.SuccessResponse{Data: orgGet, Type: appresponse.THdlSuccess, Message: appresponse.HdlMsgSuccess},
			}, Ok: true,
		}
	})
}

func (h *HandlerImpl) List() gin.HandlerFunc {
	return handler.HandlerResolver(func(c *gin.Context) handler.HandlerResponse {
		ctx := c.Request.Context()
		uFetched, err := h.OrganizationService.List(ctx)
		if err != nil {
			_, errResp, apperr := appresponse.PrepareErr(err, helper.TraceCurrentFuncArgs(map[string]string{}))
			sc := http.StatusInternalServerError

			if apperr.Type == apperror.BadRequest || apperr.Type == apperror.NotFound {
				sc = http.StatusBadRequest
				errResp = appresponse.HdlRespBadRequest()
			} else {
				sc = http.StatusInternalServerError
				errResp = appresponse.HdlRespInternalServerError()
			}

			return handler.HandlerResponse{
				ResponseWrapper: appresponse.ResponseWrapper{
					StatusCode: sc,
					Response:   errResp,
				}, Ok: false,
			}
		}

		return handler.HandlerResponse{
			ResponseWrapper: appresponse.ResponseWrapper{
				StatusCode: http.StatusOK,
				Response:   appresponse.SuccessResponse{Data: uFetched, Type: appresponse.THdlSuccess, Message: appresponse.HdlMsgSuccess},
			}, Ok: true,
		}
	})
}

// HardDelete User
// @Summary Handling request from the client to perform hard delete user
// @Tags userDelete
// @Accept json
func (h *HandlerImpl) HardDelete() gin.HandlerFunc {
	return handler.HandlerResolver(func(c *gin.Context) handler.HandlerResponse {
		orgIDParam := c.Param("organization_id")

		orgID, err := uuid.Parse(orgIDParam)
		if err != nil {
			// errMsg := fmt.Sprintf("%s: failed parse user-uuid :%v", helper.TraceCurrentFunc(), userIDParam)
			return handler.HandlerResponse{
				ResponseWrapper: appresponse.ResponseWrapper{
					StatusCode: http.StatusInternalServerError,
					Response:   appresponse.ErrorResponse{Type: appresponse.THdlBadRequest, Message: "failed parse user-uuid"},
				},
				Ok: false,
			}
		}

		ctx := c.Request.Context()
		err = h.OrganizationService.HardDelete(ctx, orgID)
		if err != nil {
			_, errResp, apperr := appresponse.PrepareErr(err, helper.TraceCurrentFuncArgs(map[string]string{}))
			sc := http.StatusInternalServerError

			if apperr.Type == apperror.BadRequest || apperr.Type == apperror.NotFound {
				sc = http.StatusBadRequest
				errResp = appresponse.HdlRespBadRequest()
			} else {
				sc = http.StatusInternalServerError
				errResp = appresponse.HdlRespInternalServerError()
			}

			return handler.HandlerResponse{
				ResponseWrapper: appresponse.ResponseWrapper{
					StatusCode: sc,
					Response:   errResp,
				}, Ok: false,
			}
		}

		return handler.HandlerResponse{
			ResponseWrapper: appresponse.ResponseWrapper{
				StatusCode: http.StatusOK,
				Response:   appresponse.SuccessResponse{Data: nil, Type: appresponse.THdlSuccess, Message: appresponse.HdlMsgSuccess},
			}, Ok: true,
		}
	})
}
