package v1

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kordape/ottct-main-service/internal/handler"
	"github.com/kordape/ottct-main-service/pkg/api"
	"github.com/kordape/ottct-main-service/pkg/httpserver"
	"github.com/kordape/ottct-main-service/pkg/token"
)

func (r *routes) updateSubscriptionsHandler(entityManager *handler.EntityManager, subscriptionsManager *handler.SubscriptionManager, tokenManager *token.Manager) func(c *gin.Context) {
	return func(c *gin.Context) {
		logger := httpserver.GetLogger(c)

		entity, err := entityManager.GetEntity(c.Param("entityid"), logger)
		if err != nil {
			logger.WithError(err).Error("Error while getting entity")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if entity == nil {
			logger.Error("Entity not found")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		var request api.UpdateSubscriptionRequest
		requestBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.WithError(err).Error("Error while reading request body")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(requestBody, &request)
		if err != nil {
			logger.WithError(err).Error("Error while unmarshaling UpdateSubscription request")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		claims, err := tokenManager.GetClaimsFromJWT(c.GetHeader("Authorization"))
		if err != nil {
			logger.WithError(err).Error("Error getting claims from bearer token")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		err = subscriptionsManager.UpdateSubscription(claims.User, entity.Id, request, logger)
		if err != nil {
			if errors.Is(err, handler.ErrInvalidRequest) {
				logger.WithError(err).Error("Error while updating subscription")
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			logger.WithError(err).Error("Error while updating subscription")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusOK)
	}
}
