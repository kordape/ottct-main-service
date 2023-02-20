package v1

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kordape/ottct-main-service/pkg/token"
)

const authHeaderKey = "Authorization"

func AuthMiddleware(tokenManager *token.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader(authHeaderKey)
		token := strings.Split(bearerToken, "Bearer ")

		if err := tokenManager.VerifyJWT(token[1]); err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
