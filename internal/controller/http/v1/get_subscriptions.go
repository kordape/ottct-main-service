package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kordape/ottct-main-service/internal/handler"
	"github.com/kordape/ottct-main-service/pkg/api"
)

func (r *routes) getSubscriptionsHandler(subscriptionsManager handler.SubscriptionManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		r.l.Debug("Get subscriptions by user request received")
		userId, err := strconv.ParseUint(c.Param("userid"), 10, 64)
		if err != nil {
			r.l.Error(fmt.Errorf("Error parsing int from userid path param: %w", err))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		entities, err := subscriptionsManager.GetSubscriptionsByUser(uint(userId))
		if err != nil {
			r.l.Error(fmt.Errorf("GetSubscriptionsByUser internal error: %w", err))
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
