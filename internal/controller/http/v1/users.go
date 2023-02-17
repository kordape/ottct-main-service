package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kordape/ottct-main-service/internal/handler"
	"github.com/kordape/ottct-main-service/pkg/api"
)

func (r *routes) newPostUsersHandler(userManager handler.UserManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		r.l.Debug("POST users request received")
		// TODO: parse request
		err := userManager.CreateUser()

		if err != nil {
			r.l.Error("POST users error: %w", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		}
		c.JSON(http.StatusOK, api.EchoResponse{
			Message: "Users Response",
		})
	}
}
