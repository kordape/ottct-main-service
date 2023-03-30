package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kordape/ottct-main-service/internal/handler"
	"github.com/kordape/ottct-main-service/pkg/api"
)

func (r *routes) newGetTweetsHandler(manager *handler.TwitterManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		r.l.Debug("GetTweets request received")

		request := api.GetTweetsRequest{}
		requestBody, _ := ioutil.ReadAll(c.Request.Body)
		r.l.Info(string(requestBody))

		err := json.Unmarshal(requestBody, &request)
		if err != nil {
			r.l.Error(fmt.Errorf("Error while unmarshaling GetTweets request: %v", err))
			c.AbortWithStatusJSON(http.StatusBadRequest, api.GetTweetsResponse{
				Error: err.Error(),
			})
			return
		}

		resp, err := manager.GetTweets(c.Request.Context(), request)

		if err != nil {
			if errors.Is(err, handler.ErrInvalidRequest) {
				r.l.Error(fmt.Errorf("Invalid GetTweets request: %v", err))
				c.AbortWithStatusJSON(http.StatusBadRequest, api.GetTweetsResponse{
					Error: err.Error(),
				})
				return
			}

			r.l.Error(fmt.Errorf("GetTweets internal error: %v", err.Error()))
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}
