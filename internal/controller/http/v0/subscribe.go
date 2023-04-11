package v0

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kordape/ottct-main-service/internal/handler"
	v1 "github.com/kordape/ottct-main-service/pkg/api"
	api "github.com/kordape/ottct-main-service/pkg/api/v0"
	"github.com/kordape/ottct-main-service/pkg/httpserver"
	"github.com/sethvargo/go-password/password"
)

func (r *routes) newPostSubscribeHandler(entityManager *handler.EntityManager, userManager *handler.AuthManager, subscriptionsManager *handler.SubscriptionManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		logger := httpserver.GetLogger(c)

		request := api.SubscribeRequest{}
		requestBody, _ := ioutil.ReadAll(c.Request.Body)

		err := json.Unmarshal(requestBody, &request)
		if err != nil {
			logger.WithError(err).Error("Error while unmarshaling SubscribeRequest")
			c.AbortWithStatusJSON(http.StatusBadRequest, api.SubscribeResponse{
				Error: err.Error(),
			})
			return
		}

		entity, err := entityManager.GetEntity(request.EntityId, logger)
		if err != nil {
			if errors.Is(err, handler.ErrEntityNotFound) {
				c.AbortWithStatusJSON(http.StatusBadRequest, api.SubscribeResponse{
					Error: err.Error(),
				})
				return
			} else {
				logger.WithError(err).Error("Error while getting entity")
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}

		if entity == nil {
			logger.Error("Entity not found")
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// check if user already exists
		var user handler.User
		user, err = userManager.GetUserByEmail(request.Email, logger)
		if err != nil {
			if errors.Is(err, handler.ErrUserNotFound) {
				// create temp user
				pswd, err := password.Generate(31, 10, 10, false, false)
				if err != nil {
					logger.WithError(err).Error("Error while generating password")
					c.AbortWithStatusJSON(http.StatusInternalServerError, api.SubscribeResponse{
						Error: err.Error(),
					})
					return
				}
				user, err = userManager.SignUp(v1.SignUpRequest{
					Email:    request.Email,
					Password: pswd,
				}, logger)
				if err != nil {
					logger.WithError(err).Error("Error while creating a temporary user")
					c.AbortWithStatusJSON(http.StatusInternalServerError, api.SubscribeResponse{
						Error: err.Error(),
					})
					return
				}

			} else {
				logger.WithError(err).Error("Error while getting a temporary user")
				c.AbortWithStatusJSON(http.StatusInternalServerError, api.SubscribeResponse{
					Error: err.Error(),
				})
			}
		}

		err = subscriptionsManager.UpdateSubscription(user.Id, request.EntityId, v1.UpdateSubscriptionRequest{
			Subscribe: true,
		}, logger)

		if err != nil {
			if errors.Is(err, handler.ErrInvalidRequest) || errors.Is(err, handler.ErrEntityNotFound) {
				logger.WithError(err).Error("Error while updating subscription")
				c.AbortWithStatusJSON(http.StatusBadRequest, api.SubscribeResponse{
					Error: err.Error(),
				})
				return
			}

			logger.WithError(err).Error("Error while updating subscription")
			c.AbortWithStatusJSON(http.StatusInternalServerError, api.SubscribeResponse{
				Error: err.Error(),
			})
			return
		}

		c.Status(http.StatusOK)
	}
}
