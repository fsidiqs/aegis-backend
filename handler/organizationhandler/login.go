package organizationhandler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fsidiqs/aegis-backend/handler"
	"github.com/fsidiqs/aegis-backend/helper"
	"github.com/fsidiqs/aegis-backend/model"
	"github.com/fsidiqs/aegis-backend/model/apperror"
	"github.com/fsidiqs/aegis-backend/model/appresponse"
)

// Login used to authenticate user
func (h *HandlerImpl) Login(c *gin.Context) {
	// swagger:operation POST /auth/login auth Login
	//
	// log in the user which eventually will return a new pair of tokens (auth and refresh token)
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	//	- name: body
	//   in: body
	//   required: true
	//   default: '{"email":"fajar@mindtera.com", "password": "12345678"}'
	//	responses:
	//   "default":
	//     description: degault response
	//     schema:
	//       type: object
	//	  "200":
	//     description: login response
	//     schema:
	//       type: object
	//       items:
	//         "$ref": "#/responses/TokenData"
	//   "401":
	//     "$ref": "#/responses/ErrorResponse"

	var req model.UserLoginReq
	// publicTokenVal interface{}
	// publicTokOK bool

	// publicTokenVal, publicTokOK = c.Get("public_user")
	// _, publicTokOK = c.Get("public_user")
	// if !publicTokOK {
	// 	errMsg := fmt.Sprintf("%s: failed get key-value: public_user", helper.TraceCurrentFunc())
	// 	h.log.Error(errMsg)

	// 	c.JSON(http.StatusInternalServerError, appresponse.ErrorResponse{Message: appresponse.HdlMsgFailedExtractUserCtx, Type: appresponse.THandlerInternal})
	// 	return
	// }
	// publicToken := publicTokenVal.(*model.ValidatedPublicToken)

	if ok := handler.BindData(c, &req); !ok {
		return
	}

	u := &model.User{
		Email:    req.Email,
		Password: req.Password,
	}

	ctx := c.Request.Context()
	// perform UserService Register with transaction
	uFetched, err := h.UserService.ComparePassword(ctx, u.Email, u.Password)
	// currently no need to use transaction since it is read only
	if err != nil {
		_, errResp, apperr := appresponse.PrepareErr(err, helper.TraceCurrentFuncArgs(req))
		sc := http.StatusInternalServerError

		switch apperr.Type {
		case apperror.BadRequest:
			// h.log.Info(internalMsg)
			sc = http.StatusUnauthorized
			errResp = appresponse.ErrorResponse{Type: appresponse.THdlUnauthorized, Message: appresponse.HdlMsgUnauthorized}
		case apperror.TUserUnverified:
			// h.log.Info(internalMsg)
			sc = http.StatusUnauthorized
			errResp = appresponse.ErrorResponse{Type: appresponse.THdlUnverified, Message: appresponse.HdlMsgUserUnverified}
		case apperror.TUserNotFound:
			// h.log.Info(internalMsg)
			sc = http.StatusUnauthorized
			errResp = appresponse.ErrorResponse{Type: appresponse.THandlerNoUser, Message: appresponse.HdlMsgNoUser}
		case apperror.NotFound:
			// h.log.Info(internalMsg)
			sc = http.StatusUnauthorized
			errResp = appresponse.ErrorResponse{Type: appresponse.THandlerNoUser, Message: appresponse.HdlMsgNoUser}
		case apperror.Authorization:
			// h.log.Info(internalMsg)
			sc = http.StatusUnauthorized
			errResp = appresponse.HdlRespUnautthorized()
		default:

			// h.log.Error(internalMsg)
			sc = http.StatusInternalServerError
			errResp = appresponse.HdlRespInternalServerError()
		}

		c.JSON(sc, errResp)
		return

	}

	// check if session available
	// if available, then just update
	// uSess, err := h.UserService.FindUserSessionDBByUserID(ctx, uFetched.ID)
	// if err != nil {
	// 	log.Printf("Failed to sign in user:Email:%v %v\n", uFetched.Email, err)
	// 	c.JSON(apperror.Status(err), appresponse.ErrorResponse{Message: appresponse.HdlMsgInternal, Type: appresponse.THandlerInternal})
	// 	return
	// }

	tokens, err := h.TokenService.NewPairFromUser(ctx, uFetched)
	if err != nil {
		log.Printf("Failed to create tokens for user: %v\n", err.Error())
		c.JSON(apperror.Status(err), appresponse.ErrorResponse{Message: appresponse.HdlMsgInternal, Type: appresponse.THandlerInternal})
		return
	}

	// session exists, update it
	// redis: remove and set new session
	// db: update
	// if uSess != nil {
	// 	err := trx.RedisRemoveSession(ctx, uSess.UserID)
	// 	if err != nil {
	// 		trx.RollbackTrx()
	// 		log.Printf("Failed to sign in user: %v\n", err)
	// 		c.JSON(apperror.Status(err), appresponse.ErrorResponse{Message: appresponse.HdlMsgInternal, Type: appresponse.THandlerInternal})
	// 		return
	// 	}
	// 	now := time.Now()
	// 	uSessUpdate := model.UserSessionUpdate{
	// 		NotificationToken: publicToken.NotificationToken,
	// 		AuthToken:         tokens.AuthToken.SS,
	// 		RefreshTokenID:    tokens.RefreshToken.ID.String(),
	// 		ExpiredAt:         now.Add(tokens.ExpiresIn),
	// 		Status:            model.TUserSessionActive,
	// 		AccountType:       uFetched.AccountType,
	// 		UserRecordFlag:    uFetched.RecordFlag,

	// 		DefaultColumns: model.DefaultColumns{
	// 			UpdatedBy: uFetched.ID.String(),
	// 			UpdatedAt: &now,
	// 		},
	// 	}
	// 	err = trx.UpdateUserSession(ctx, uSess.UserID, uSessUpdate)
	// 	if err != nil {
	// 		trx.RollbackTrx()
	// 		log.Printf("Failed to update user session:userid:%v %v\n", uSess.UserID, err)
	// 		c.JSON(apperror.Status(err), appresponse.ErrorResponse{Message: appresponse.HdlMsgInternal, Type: appresponse.THandlerInternal})
	// 		return
	// 	}
	// } else {
	// session doesnt exist, create a new session

	// store session to db and redis

	// uSess := model.NewUserSession(tokens,
	// 	model.ExtraData{
	// 		NotificationToken: publicToken.NotificationToken,
	// 		Status:            model.TUserSessionActive,
	// 		AccountType:       uFetched.AccountType,
	// 		UserRecordFlag:    uFetched.RecordFlag,
	// 	})
	// TODO? might be improved use goroutine and send through channel to userservice.register
	// err = trx.StoreSession(ctx, &uSess)
	// if err != nil {
	// 	trx.RollbackTrx()

	// 	log.Printf("Failed to store token to db and redis: %v\n", err.Error())

	// 	// may eventually implement rollback logic here
	// 	// meaning, if we fail to create tokens after creating a user,
	// 	// we make sure to clear/delete the created user in the databse

	// 	c.JSON(apperror.Status(err), appresponse.ErrorResponse{Message: appresponse.HdlMsgInternal, Type: appresponse.THandlerInternal})
	// 	return
	// }

	// }
	// lastLoginMethod := &model.UserUpdate{
	// 	LastLoginMethod: string(model.TLastLoginManual),
	// 	DefaultColumns: model.DefaultColumns{
	// 		UpdatedBy: helper.TraceCurrentFunc(),
	// 	},
	// }
	// err = trx.UpdateUser(ctx, uFetched.ID, lastLoginMethod)
	// if err != nil {
	// 	trx.RollbackTrx()

	// 	log.Printf("Failed UpdateUser:email:%v %v\n", req.Email, err.Error())

	// 	c.JSON(apperror.Status(err), appresponse.ErrorResponse{Message: appresponse.HdlMsgInternal, Type: appresponse.THandlerInternal})
	// 	return
	// }

	// // needs to commit to avoid not found user id when creating activity log
	// errCommit := trx.CommitTrx()

	// if errCommit != nil {
	// 	log.Printf("Failed to login:email%v %v\n", u.Email, err)

	// 	c.JSON(apperror.Status(err), appresponse.ErrorResponse{Message: appresponse.HdlMsgInternal, Type: appresponse.THandlerInternal})
	// 	return
	// }
	// // user_activity_log

	// err = h.ActivityLogService.Store(ctx, model.UserActivityLog{
	// 	UserId:            uFetched.ID,
	// 	ActivityDate:      time.Now(),
	// 	ActivityType:      "users",
	// 	RelatedActivityID: nil,
	// 	DefaultColumns:    model.NewLogDefaultCol(helper.TraceCurrentFunc()),
	// })
	// if err != nil {
	// 	log.Printf("Failed to login: %v\n", err)

	// 	c.JSON(apperror.Status(err), appresponse.ErrorResponse{Message: appresponse.HdlMsgInternal, Type: appresponse.THandlerInternal})
	// 	return
	// }
	// append notification to the response
	// tokens.NotificationToken = publicToken.NotificationToken
	c.JSON(http.StatusOK, appresponse.SuccessResponse{Data: tokens})
}
