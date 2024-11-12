package candhisapi

import (
	"github.com/gin-gonic/gin"
	"github.com/tul1/candhis_api/internal/application/repository"
	"github.com/tul1/candhis_api/openapi"
)

type candhisAPI struct {
	router   *gin.Engine
	waveData repository.WaveData
}

func NewCandhisAPI(e *gin.Engine, waveData repository.WaveData) *candhisAPI {
	api := candhisAPI{router: e, waveData: waveData}
	openapi.RegisterHandlers(e, api)
	return &api
}
