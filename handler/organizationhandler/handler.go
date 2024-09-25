package organizationhandler

import (
	"github.com/gin-gonic/gin"

	"github.com/fsidiqs/aegis-backend/handler/middleware"
	"github.com/fsidiqs/aegis-backend/service"
)

type HandlerImpl struct {
	OrganizationService service.IOrganizationService
	UserService         service.IUserService
	TokenService        service.ITokenService
}

type Config struct {
	R                   *gin.Engine
	OrganizationService service.IOrganizationService
	UserService         service.IUserService
	TokenService        service.ITokenService
	BaseURL             string
}

func NewHandler(c *Config) {
	h := &HandlerImpl{
		OrganizationService: c.OrganizationService,
		TokenService:        c.TokenService,
		UserService:         c.UserService,
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

	g.POST("/organization", authmiddleware.AuthUser(), h.Create())
	// g.GET("/organization", authmiddleware.AuthUser(), h.get())
	g.GET("/organizations", authmiddleware.AuthUser(), h.List())
	// g.POST("/user/resend-verification", h.ResendEmailVerification)

	g.GET("/organization/:organization_id", authmiddleware.AuthUser(), h.Get())
	g.PUT("/organization/:organization_id", authmiddleware.AuthUser(), h.Update())
	g.DELETE("/organization/:organization_id", authmiddleware.AuthUser(), h.HardDelete())
}
