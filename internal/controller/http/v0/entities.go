package v0

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kordape/ottct-main-service/internal/handler"
	"github.com/kordape/ottct-main-service/pkg/api"
	"github.com/kordape/ottct-main-service/pkg/httpserver"
)

func (r *routes) getEntitiesHandler(entityManager *handler.EntityManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		logger := httpserver.GetLogger(c)

		entities, err := entityManager.GetEntities(logger)
		if err != nil {
			logger.WithError(err).Error("GetEntities internal error")
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}

		response := make(api.EntitiesResponse, len(entities))
		for i, e := range entities {
			response[i] = api.Entity{
				Id:            e.Id,
				TwitterId:     e.TwitterId,
				DisplayName:   e.DisplayName,
				TwitterHandle: e.Handle,
			}
		}

		c.JSON(http.StatusOK, response)
	}
}
