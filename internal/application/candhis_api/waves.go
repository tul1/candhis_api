package candhisapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tul1/candhis_api/openapi"
)

func (s candhisAPI) Waves(c *gin.Context) {
	errorResponse := openapi.Error{Message: "unimplemented path"}
	c.JSON(http.StatusInternalServerError, errorResponse)
}
