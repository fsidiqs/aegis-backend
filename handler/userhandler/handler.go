package userhandler

import (
	"github.com/gin-gonic/gin"

	"github.com/fsidiqs/aegis-backend/handler/middleware"
	"github.com/fsidiqs/aegis-backend/mail"
	"github.com/fsidiqs/aegis-backend/service"
	"github.com/fsidiqs/aegis-backend/service/tokenservice"
)

type HandlerImpl struct {
	PublicTokenService tokenservice.IPublicTokenService
	UserService        service.IUserService
	TokenService       service.ITokenService
	MailClient         mail.IMailClient
}

type Config struct {
	R                  *gin.Engine
	UserService        service.IUserService
	TokenService       service.ITokenService
	PublicTokenService tokenservice.IPublicTokenService
	MailClient         mail.IMailClient
	BaseURL            string
}

func NewHandler(c *Config) {
	h := &HandlerImpl{
		UserService:        c.UserService,
		TokenService:       c.TokenService,
		PublicTokenService: c.PublicTokenService,
		MailClient:         c.MailClient,
	}

	// make auth middleware
	authmiddleware := middleware.AuthMiddleware{
		TokenService: c.TokenService,
		UserService:  c.UserService,
	}

	// authPublicMid := middleware.AuthPublicMiddleware{
	// 	PublicTokenService: c.PublicTokenService,
	// }

	g := c.R.Group(c.BaseURL)

	g.POST("/auth/login", h.Login)
	// g.POST("/auth/logout", authmiddleware.AuthUser(), h.Logout)

	g.POST("/auth/tokens", h.Tokens)

	// swagger:operation POST /user user Register
	// ---
	//	summary: Register a new user and send them an email verification expires in 1 day
	// consumes:
	// - application/json
	// produces:
	// - application/json
	//	parameters:
	//	- name: name
	//	  in: body
	//   description: abc
	//	  type: string
	//	  required: true
	//	- name: email
	//	  in: body
	//   description: abc
	//	  type: string
	//	  required: true
	//	- name: phone_number
	//	  in: body
	//   description: abc
	//	  type: string
	//	  required: true
	//	- name: password
	//	  in: body
	//   description: abc
	//	  type: string
	//	  required: true
	// responses:
	//   "201":
	//   "401":
	//     "$ref": "#/responses/ErrorResponse"
	g.POST("/user", h.Register())
	g.PUT("/user", authmiddleware.AuthUser(), h.UpdateMyDetails())
	g.GET("/user", authmiddleware.AuthUser(), h.MyDetails())
	g.GET("/users", authmiddleware.AuthUser(), h.ListUsers())
	// g.POST("/user/resend-verification", h.ResendEmailVerification)

	g.GET("/user/:user_id", authmiddleware.AuthUser(), h.Details())
	g.PUT("/user/:user_id", authmiddleware.AuthUser(), h.UpdateDetails())
	g.DELETE("/user/:user_id", authmiddleware.AuthUser(), h.HardDeleteUser())

	// g.POST("/user/verify-email-by-otp", authPublicMid.AuthPublic(), h.VerifyEmailByOTP)
	g.POST("/user/forgot-password", h.ForgotPasswordByEmail)
	g.POST("/user/forgot-password/otp", h.SubmitOTPForgotPassword)
	g.POST("/user/forgot-password/update-using-otp", h.UpdateForgottenPasswordByOTP)
}
