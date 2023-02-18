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

func (r *routes) newSignUpHandler(userManager handler.AuthManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		r.l.Debug("SignUp request received")

		request := api.SignUpRequest{}
		requestBody, _ := ioutil.ReadAll(c.Request.Body)

		err := json.Unmarshal(requestBody, &request)
		if err != nil {
			r.l.Error(fmt.Errorf("Error while unmarshaling SignUp request: %v", err))
			c.AbortWithStatusJSON(http.StatusBadRequest, api.SignUpResponse{
				Error: err.Error(),
			})
			return
		}

		err = userManager.SignUp(request)
		if errors.Is(err, handler.ErrInvalidRequest) {
			r.l.Error(fmt.Errorf("Invalid SignUp request: %v", err))
			c.AbortWithStatusJSON(http.StatusBadRequest, api.SignUpResponse{
				Error: err.Error(),
			})
			return
		}

		if err != nil {
			r.l.Error(fmt.Errorf("SignUp internal error: %v", err.Error()))
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		}
		c.JSON(http.StatusOK, api.SignUpResponse{
			Message: "SignUp Response",
		})
	}
}
