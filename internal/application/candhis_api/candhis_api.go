package candhisapi

import (
	"github.com/gin-gonic/gin"
	"github.com/tul1/candhis_api/openapi"
)

type candhisAPI struct {
	router *gin.Engine
}

func NewCandhisAPI(e *gin.Engine) *candhisAPI {
	api := candhisAPI{router: e}
	openapi.RegisterHandlers(e, api)
	return &api
}
