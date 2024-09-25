package organizationhandler

// func (h *HandlerImpl) VerifyEmailByOTP(c *gin.Context) {
// 	var (
// 		req            model.ReqVerifyEmailByOTP
// 		publicTokenVal interface{}
// 		publicTokOK    bool
// 	)

// 	publicTokenVal, publicTokOK = c.Get("public_user")
// 	if !publicTokOK {
// 		log.Printf("failed to get public_user")

// 		c.JSON(http.StatusInternalServerError, appresponse.ErrorResponse{Message: appresponse.HdlMsgFailedExtractUserCtx, Type: appresponse.THandlerInternal})
// 		return
// 	}
// 	publicToken := publicTokenVal.(*model.ValidatedPublicToken)

// 	if ok := handler.BindData(c, &req); !ok {
// 		return
// 	}

// 	ctx := c.Request.Context()
// 	// perform UserService with transaction
// 	uFetch, err := h.UserService.VerifyEmail(ctx, req.Email, req.OTP)
// 	if err != nil {
// 		trx.RollbackTrx()

// 		internalMsg, errResp, apperr := appresponse.PrepareErr(err, helper.TraceCurrentFuncArgs(req))

// 		sc := http.StatusInternalServerError
// 		if apperr.Type == apperror.BadRequest {
// 			h.log.Info(internalMsg)
// 			sc = http.StatusBadRequest
// 			errResp = appresponse.HdlRespBadRequest()
// 		} else if apperr.Type == apperror.EmailNotFound {
// 			sc = http.StatusBadRequest
// 			errResp = appresponse.ErrorResponse{Message: appresponse.MsgEmailNotFound, Type: appresponse.TEmailNotFound}
// 		} else if apperr.Type == apperror.OTPNotFound {
// 			sc = http.StatusBadRequest
// 			errResp = appresponse.ErrorResponse{Message: appresponse.MsgOTPNotFound, Type: appresponse.TOTPNotFound}
// 		} else if apperr.Type == apperror.OTPExpired {
// 			sc = http.StatusBadRequest
// 			errResp = appresponse.ErrorResponse{Message: appresponse.HdlMsgOTPExpire, Type: appresponse.TOTPExpired}
// 		} else {
// 			h.log.Error(internalMsg)
// 			sc = http.StatusInternalServerError
// 			errResp = appresponse.HdlRespInternalServerError()
// 		}

// 		c.JSON(sc, errResp)
// 		return

// 	}

// 	tokens, err := h.TokenService.NewPairFromUser(ctx, uFetch)
// 	if err != nil {
// 		trx.RollbackTrx()
// 		log.Printf("Failed to create tokens for user: %v\n", err.Error())
// 		var errmsg appresponse.ErrMsgResp = appresponse.ErrMsgResp(err.Error())
// 		c.JSON(apperror.Status(err), appresponse.ErrorResponse{Message: errmsg})
// 		return
// 	}

// 	uSess := model.NewUserSession(tokens,
// 		model.ExtraData{
// 			NotificationToken: publicToken.NotificationToken,
// 			Status:            model.TUserSessionActive,
// 			AccountType:       uFetch.AccountType,
// 			UserRecordFlag:    uFetch.RecordFlag,
// 		})

// 	// TODO? might be improved use goroutine and send through channel to userservice.register
// 	err = trx.StoreSession(ctx, &uSess)
// 	if err != nil {
// 		trx.RollbackTrx()

// 		log.Printf("Failed StoreSession:email:%v %v\n", req.Email, err.Error())

// 		// may eventually implement rollback logic here
// 		// meaning, if we fail to create tokens after creating a user,
// 		// we make sure to clear/delete the created user in the databse

// 		var errmsg appresponse.ErrMsgResp = appresponse.ErrMsgResp(err.Error())
// 		c.JSON(apperror.Status(err), appresponse.ErrorResponse{Message: errmsg})
// 		return
// 	}

// 	lastLoginMethod := &model.UserUpdate{
// 		LastLoginMethod: string(model.TLastLoginManual),
// 		DefaultColumns: model.DefaultColumns{
// 			UpdatedBy: helper.TraceCurrentFunc(),
// 		},
// 	}
// 	err = trx.UpdateUser(ctx, uFetch.ID, lastLoginMethod)
// 	if err != nil {
// 		trx.RollbackTrx()

// 		log.Printf("Failed UpdateUser:email:%v %v\n", req.Email, err.Error())

// 		var errmsg appresponse.ErrMsgResp = appresponse.ErrMsgResp(err.Error())
// 		c.JSON(apperror.Status(err), appresponse.ErrorResponse{Message: errmsg})
// 		return
// 	}

// 	// user_activity_log

// 	err = h.ActivityLogService.Store(ctx, model.UserActivityLog{
// 		UserId:            uFetch.ID,
// 		ActivityDate:      time.Now(),
// 		ActivityType:      "users",
// 		RelatedActivityID: nil,
// 		DefaultColumns:    model.NewLogDefaultCol(helper.TraceCurrentFunc()),
// 	})
// 	if err != nil {
// 		trx.RollbackTrx()

// 		log.Printf("Failed to store token to db and redis: %v\n", err.Error())

// 		var errmsg appresponse.ErrMsgResp = appresponse.ErrMsgResp(err.Error())
// 		c.JSON(apperror.Status(err), appresponse.ErrorResponse{Message: errmsg})
// 		return
// 	}

// 	errCommit := trx.CommitTrx()

// 	if errCommit != nil {
// 		log.Printf("Failed to resend user verify email: %v\n", err)

// 		var errmsg appresponse.ErrMsgResp = appresponse.ErrMsgResp(err.Error())
// 		c.JSON(apperror.Status(err), appresponse.ErrorResponse{Message: errmsg})
// 		return
// 	}
// 	// append notification token to the response
// 	tokens.NotificationToken = publicToken.NotificationToken

// 	c.JSON(http.StatusOK, appresponse.SuccessResponse{Data: tokens})
// }
