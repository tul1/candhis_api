package candhisapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tul1/candhis_api/openapi"
)

func (s candhisAPI) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, openapi.Pong{Message: "pong"})
}
