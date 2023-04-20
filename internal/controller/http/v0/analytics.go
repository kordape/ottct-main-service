package v0

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kordape/ottct-main-service/internal/handler"
	"github.com/kordape/ottct-main-service/pkg/api"
	"github.com/kordape/ottct-main-service/pkg/httpserver"
)

func (r *routes) newGetAnalyticsHandler(manager *handler.TwitterManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		logger := httpserver.GetLogger(c)

		request := api.GetTweetsRequest{
			EntityID:   c.Query("entityId"),
			From:       toTimeOrZero(c.Query("from")),
			To:         toTimeOrZero(c.Query("to")),
			MaxResults: toIntOrZero(c.Query("maxResults")),
		}

		resp, err := manager.GetTweets(c.Request.Context(), request, logger)

		if err != nil {
			if errors.Is(err, handler.ErrInvalidRequest) {
				logger.WithError(err).Error("Invalid GetTweets request")
				c.AbortWithStatusJSON(http.StatusBadRequest, api.GetTweetsResponse{
					Error: err.Error(),
				})
				return
			}

			logger.WithError(err).Error("GetTweets internal error")
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

func toTimeOrZero(in string) time.Time {
	var err error
	out := time.Time{}

	if in != "" {
		out, err = time.Parse(time.RFC3339, in)
		if err != nil {
			return out
		}
	}

	return out
}

func toIntOrZero(in string) int {
	out, err := strconv.Atoi(in)
	if err != nil {
		return 0
	}

	return out
}
