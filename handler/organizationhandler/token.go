package organizationhandler

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/fsidiqs/aegis-backend/handler"
	"github.com/fsidiqs/aegis-backend/helper"
	"github.com/fsidiqs/aegis-backend/model"
	"github.com/fsidiqs/aegis-backend/model/apperror"
	"github.com/fsidiqs/aegis-backend/model/appresponse"
)

// Tokens handler

func (h *HandlerImpl) Tokens(c *gin.Context) {
	// swagger:operation POST /auth/tokens auth Tokens
	//
	// refreshes current tokens (auth and refresh token)
	//
	// 1. validates `refresh_token` in request body, if valid
	// 2. invoke user service, to get user from current RefreshToken
	// 3. if the user is unverified or does not have ACTIVE flag then return 401 error code response
	// 4. check current token in redis and db, if doesn't exist in redis we keep looking the user session in the database. if token is not exists in both redis and db then return 401 error code response
	// 5. remove the current session from redis and db
	// 6. create a new pair token (auth_token and refresh_token)
	// 7. create a new user session, store it in redis and db
	// 8. if the current operation succeed, then returns a new pair tokens (auth_token and refresh_token)
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	//	parameters:
	//	- name: refresh_token
	//	  in: body
	//	  type: string
	//	  required: true
	//	responses:
	//	  "200":
	//     "$ref": "#/responses/TokenData"
	//	  "401":
	//     "$ref": "#/responses/ErrorResponse"

	// bind JSON to req of type tokensRew
	var req model.TokensReq

	if ok := handler.BindData(c, &req); !ok {
		return
	}

	ctx := c.Request.Context()
	// perform UserService Register with transaction

	// verify refresh JWT
	refreshToken, err := h.TokenService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {

		_, errResp, apperr := appresponse.PrepareErr(err, helper.TraceCurrentFuncArgs(req))
		sc := http.StatusInternalServerError

		if apperr.Type == apperror.BadRequest || apperr.Type == apperror.NotFound || apperr.Type == apperror.Authorization {
			// h.log.Info(internalMsg)
			sc = http.StatusUnauthorized
			errResp = appresponse.HdlRespBadRequest()
		} else {
			// h.log.Error(internalMsg)
			sc = http.StatusInternalServerError
			errResp = appresponse.HdlRespInternalServerError()
		}

		c.JSON(sc, errResp)
		return
	}

	// get up-to-date user
	uFetched, err := h.UserService.Get(ctx, refreshToken.UID)
	// check if user is verified
	if err != nil {

		_, errResp, apperr := appresponse.PrepareErr(err, helper.TraceCurrentFuncArgs(req))
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
	if uFetched.EmailVerifiedAt == nil || uFetched.RecordFlag != model.RecActive {

		c.JSON(http.StatusUnauthorized, appresponse.ErrorResponse{Message: "user is either unverified or not ACTIVE"})
		return
	}

	// exists, err := h.UserService.RefreshTokenExists(ctx, uFetched.ID, refreshToken.ID)
	// if err != nil {

	// 	log.Printf("Failed to check RefreshToken:%v, userID: %v, err: %v\n", refreshToken.ID, uFetched.ID, err.Error())

	// 	c.JSON(http.StatusUnauthorized, appresponse.ErrorResponse{Message: "token does not exist"})
	// 	return
	// }
	// if !exists {

	// 	c.JSON(http.StatusUnauthorized, appresponse.ErrorResponse{Message: "token does not exist"})
	// 	return
	// }

	tokens, err := h.TokenService.NewPairFromUser(ctx, uFetched)
	if err != nil {

		log.Printf("Failed to create tokens for user: %+v. Error: %v\n", uFetched.ID, err.Error())

		_, errResp, apperr := appresponse.PrepareErr(err, helper.TraceCurrentFuncArgs(req))
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

	now := time.Now()
	uSessUpdate := model.UserSessionUpdate{
		AuthToken:      tokens.AuthToken.SS,
		RefreshTokenID: tokens.RefreshToken.ID.String(),
		ExpiredAt:      now.Add(tokens.ExpiresIn),
		Status:         model.TUserSessionActive,
		// AccountType:    uFetched.AccountType,
		UserRecordFlag: uFetched.RecordFlag,
		DefaultColumns: model.DefaultColumns{
			UpdatedBy: uFetched.ID.String(),
			UpdatedAt: &now,
		},
	} // store session to db(update) and redis

	// TODO? might be improved use goroutine and send through channel to userservice.register
	err = h.UserService.UpdateUserSession(ctx, uFetched.ID, uSessUpdate)
	if err != nil {

		log.Printf("Failed to store token to db and redis: %v\n", err.Error())

		_, errResp, apperr := appresponse.PrepareErr(err, helper.TraceCurrentFuncArgs(req))
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

	// errCommit := h.UserService.CommitTrx()

	// if errCommit != nil {
	// 	log.Printf("Failed to resend user forgot password: %v\n", err)

	// 	_, errResp, apperr := appresponse.PrepareErr(err, helper.TraceCurrentFuncArgs(req))
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
