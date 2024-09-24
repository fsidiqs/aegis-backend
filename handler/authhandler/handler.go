package authhandler

import (
	"github.com/gin-gonic/gin"

	"github.com/fsidiqs/aegis-backend/service/tokenservice"
)

type HandlerImpl struct {
	PubTokenService tokenservice.IPublicTokenService
}

type Config struct {
	R               *gin.Engine
	PubTokenService tokenservice.IPublicTokenService
	BaseURL         string
}

func NewHandler(c *Config) error {
	h := &HandlerImpl{
		PubTokenService: c.PubTokenService,
	}

	g := c.R.Group(c.BaseURL)

	g.POST("/auth/public", h.CreatePublicToken)

	return nil
}
