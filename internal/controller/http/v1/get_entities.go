package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kordape/ottct-main-service/internal/handler"
	"github.com/kordape/ottct-main-service/pkg/api"
)

func (r *routes) getEntitiesHandler(entityManager handler.EntityManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		r.l.Debug("Get entities request received")

		entities, err := entityManager.GetEntities()
		if err != nil {
			r.l.Error(fmt.Errorf("GetEntities internal error: %w", err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}

		// TODO posto je handler.Entity : api.Entity mapiranje 1:1 da li da radimo ovaj blok?
		// Nepotrebno je u ovom slucaju.
		response := make(api.EntitiesResponse, len(entities))
		for i, e := range entities {
			response[i] = api.Entity{
				Id:          e.Id,
				TwitterId:   e.TwitterId,
				DisplayName: e.DisplayName,
			}
		}

		c.JSON(http.StatusOK, response)
	}
}
