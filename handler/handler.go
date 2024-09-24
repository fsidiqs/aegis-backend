package handler

import (
	"context"

	"github.com/fsidiqs/aegis-backend/model/appresponse"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HandlerResponse struct {
	Ctx context.Context
	appresponse.ResponseWrapper
	TrxKeys []string
	Ok      bool
}

type Handler func(*gin.Context) HandlerResponse

func HandlerResolver(handler Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var handlerResponse HandlerResponse
		handlerResponse = handler(c)

		if !handlerResponse.Ok {
			c.JSON(handlerResponse.StatusCode, handlerResponse.Response)
			for _, trxKey := range handlerResponse.TrxKeys {
				trx := handlerResponse.Ctx.Value(trxKey)
				if trx == nil {
					continue
				}
				trx.(*gorm.DB).Rollback()
			}
			return
		}

		for _, trxKey := range handlerResponse.TrxKeys {
			trx := handlerResponse.Ctx.Value(trxKey)
			if trx == nil {
				continue
			}
			trx.(*gorm.DB).Commit()
		}

		c.JSON(handlerResponse.StatusCode, handlerResponse.Response)
	}
}
