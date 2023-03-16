package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kordape/ottct-main-service/internal/handler"
	"github.com/kordape/ottct-main-service/pkg/api"
	"github.com/kordape/ottct-main-service/pkg/token"
)

func (r *routes) getSubscriptionsHandler(subscriptionsManager *handler.SubscriptionManager, tokenManager *token.Manager) func(c *gin.Context) {
	return func(c *gin.Context) {
		r.l.Debug("Get subscriptions by user request received")
		claims, err := tokenManager.GetClaimsFromJWT(c.GetHeader("Authorization"))
		if err != nil {
			r.l.Error(fmt.Errorf("Error getting claims from bearer token: %w", err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		entities, err := subscriptionsManager.GetSubscriptionsByUser(claims.User)
		if err != nil {
			r.l.Error(fmt.Errorf("GetSubscriptionsByUser internal error: %w", err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}

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
