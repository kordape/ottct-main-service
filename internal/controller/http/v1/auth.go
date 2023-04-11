package v1

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kordape/ottct-main-service/internal/handler"
	"github.com/kordape/ottct-main-service/pkg/api"
	"github.com/kordape/ottct-main-service/pkg/httpserver"
)

func (r *routes) newSignUpHandler(userManager *handler.AuthManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		logger := httpserver.GetLogger(c)

		request := api.SignUpRequest{}
		requestBody, _ := ioutil.ReadAll(c.Request.Body)

		err := json.Unmarshal(requestBody, &request)
		if err != nil {
			logger.WithError(err).Error("Error while unmarshaling SignUp request")
			c.AbortWithStatusJSON(http.StatusBadRequest, api.SignUpResponse{
				Error: err.Error(),
			})
			return
		}

		err = userManager.SignUp(request, logger)

		if err != nil {
			if errors.Is(err, handler.ErrInvalidRequest) {
				logger.WithError(err).Error("Invalid SignUp request")
				c.AbortWithStatusJSON(http.StatusBadRequest, api.SignUpResponse{
					Error: err.Error(),
				})
				return
			}

			logger.WithError(err).Error("SignUp internal error")
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, api.SignUpResponse{
			Message: "Success",
		})
	}
}

func (r *routes) newAuthHandler(userManager *handler.AuthManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		logger := httpserver.GetLogger(c)

		request := api.AuthRequest{}
		requestBody, _ := ioutil.ReadAll(c.Request.Body)

		err := json.Unmarshal(requestBody, &request)
		if err != nil {
			logger.WithError(err).Error("Error while unmarshaling Auth request")
			c.AbortWithStatusJSON(http.StatusBadRequest, api.AuthResponse{
				Error: err.Error(),
			})
			return
		}

		token, err := userManager.Auth(request, logger)

		if err != nil {
			if errors.Is(err, handler.ErrInvalidRequest) {
				logger.WithError(err).Error("Invalid Auth request")
				c.AbortWithStatusJSON(http.StatusBadRequest, api.AuthResponse{
					Error: err.Error(),
				})
				return
			}

			if errors.Is(err, handler.ErrUserNotFound) {
				logger.WithError(err).Error("User unauthorized")
				c.AbortWithStatusJSON(http.StatusUnauthorized, api.AuthResponse{
					Error: err.Error(),
				})
				return
			}

			logger.WithError(err).Error("Auth internal error")
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, api.AuthResponse{
			Token: token,
		})
	}
}
