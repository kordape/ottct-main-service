package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kordape/ottct-main-service/internal/handler"
	"github.com/kordape/ottct-main-service/pkg/api"
	"github.com/kordape/ottct-main-service/pkg/token"
)

func (r *routes) updateSubscriptionsHandler(entityManager *handler.EntityManager, subscriptionsManager *handler.SubscriptionManager, tokenManager *token.Manager) func(c *gin.Context) {
	return func(c *gin.Context) {
		r.l.Debug("Update subscriptions request received")

		entity, err := entityManager.GetEntity(c.Param("entityid"))
		if err != nil {
			r.l.Error(fmt.Errorf("error while getting entity: %s", err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if entity == nil {
			r.l.Info("Entity not found")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		var request api.UpdateSubscriptionRequest
		requestBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			r.l.Error(fmt.Errorf("error while reading request body: %s", err))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(requestBody, &request)
		if err != nil {
			r.l.Error(fmt.Errorf("error while unmarshaling UpdateSubscription request: %v", err))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		claims, err := tokenManager.GetClaimsFromJWT(c.GetHeader("Authorization"))
		if err != nil {
			r.l.Error(fmt.Errorf("error getting claims from bearer token: %w", err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		err = subscriptionsManager.UpdateSubscription(claims.User, entity.Id, request)
		if err != nil {
			if errors.Is(err, handler.ErrInvalidRequest) {
				r.l.Error(fmt.Errorf("error while updating subscription: %s", err))
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
			
			r.l.Error(fmt.Errorf("error while updating subscription: %s", err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusOK)
	}
}
