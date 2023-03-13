package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kordape/ottct-main-service/pkg/api"
)

func (r *routes) getEntitiesHandler(c *gin.Context) {
	r.l.Debug("Get entities request received")

	result, err := r.db.GetTwitterEntities()
	if err != nil {
		c.Error(fmt.Errorf("error getting entities from db: %w", err))
	}

	response := make(api.GetEntitiesResponse, len(result))
	for i, e := range result {
		response[i] = api.Entity{
			Id:               e.ID,
			TwitterAccountId: e.TwitterAccountId,
			Name:             e.DisplayName,
		}
	}

	c.JSON(http.StatusOK, response)
}
