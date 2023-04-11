package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kordape/ottct-main-service/internal/handler"
	"github.com/kordape/ottct-main-service/pkg/api"
	"github.com/kordape/ottct-main-service/pkg/token"
)

func (r *routes) getSubscriptionsHandler(subscriptionsManager *handler.SubscriptionManager, tokenManager *token.Manager) func(c *gin.Context) {
	return func(c *gin.Context) {
		logger := getLogger(c)
		claims, err := tokenManager.GetClaimsFromJWT(c.GetHeader("Authorization"))
		if err != nil {
			logger.WithError(err).Error("Error getting claims from bearer token")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		entities, err := subscriptionsManager.GetSubscriptionsByUser(claims.User, logger)
		if err != nil {
			logger.WithError(err).Error("GetSubscriptionsByUser internal error")
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
