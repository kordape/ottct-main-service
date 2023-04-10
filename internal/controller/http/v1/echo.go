package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kordape/ottct-main-service/pkg/api"
)

func (r *routes) echoHandler(c *gin.Context) {
	logger := getLogger(c)
	logger.Debug("Echo request received")
	c.JSON(http.StatusOK, api.EchoResponse{
		Message: "Echo Response",
	})
}
