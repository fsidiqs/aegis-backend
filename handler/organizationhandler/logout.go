package organizationhandler

// func (h *HandlerImpl) Logout(c *gin.Context) {
// 	// swagger:operation POST /auth/logout auth Logout
// 	//
// 	//
// 	// ---
// 	// consumes:
// 	// - application/json
// 	// produces:
// 	// - application/json
// 	// parameters:
// 	// - name: authorization
// 	//   description: an auth token
// 	//   in: header
// 	//   type: string
// 	//	responses:
// 	//	  "200":
// 	//     "$ref": "#/responses/TokenData"
// 	//	  "401":
// 	//     "$ref": "#/responses/ErrorResponse"

// 	user, ok := c.Get("user")

// 	if !ok {
// 		c.JSON(http.StatusInternalServerError, appresponse.ErrorResponse{Message: "failed to extract user"})
// 		return
// 	}

// 	ctx := c.Request.Context()
// 	if err := h.UserService.Logout(ctx, user.(*model.User).ID); err != nil {
// 		c.JSON(apperror.Status(err), appresponse.ErrorResponse{Message: appresponse.HdlMsgInternal})
// 		return
// 	}

// 	c.JSON(http.StatusOK, appresponse.SuccessResponse{Message: "user signed out successfully"})
// }
