package v1

import (
	"encoding/json"
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
			r.l.Error(fmt.Errorf("Error while getting entity: %s", err))
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
			r.l.Error(fmt.Errorf("Error while reading request body: %s", err))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(requestBody, &request)
		if err != nil {
			r.l.Error(fmt.Errorf("Error while unmarshaling UpdateSubscription request: %v", err))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		claims, err := tokenManager.GetClaimsFromJWT(c.GetHeader("Authorization"))
		if err != nil {
			r.l.Error(fmt.Errorf("Error getting claims from bearer token: %w", err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if request.Subscribe {
			err = subscriptionsManager.AddSubscription(claims.User, entity.Id)
		} else {
			err = subscriptionsManager.DeleteSubscription(claims.User, entity.Id)
		}

		if err != nil {
			r.l.Error(fmt.Errorf("Error while updating subscription: %s", err))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusOK)
		return
	}
}
